/*
 * JLevelDB Encrypted Storage
 *
 *    Copyright (C) 2021 Jeffrey H. Johnson <trnsz@pobox.com>
 *    Copyright (C) 2019 Tenta, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file contains some code originally from GoLevelDB licensed under
 * BSD license which bears the following copyright
 *
 * Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
 *
 * All rights reserved.
 *
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/johnsonjh/jleveldb/leveldb/storage"
)

const additionalDataLen = 1 + binary.MaxVarintLen64

var (
	errFileOpen = errors.New(
		"leveldb/storage: file still open")
	errReadOnly = errors.New(
		"leveldb/storage: storage is read only")
	errCorruptedCurrent = errors.New(
		"leveldb/storage: corrupted or incomplete CURRENT file")
	errNonceUnavailable = errors.New(
		"leveldb/aesgcm: unable to generate a nonce")
)

type aesgcmStorage struct {
	path     string
	readOnly bool

	mu    sync.Mutex
	flock fileLock
	slock *aesgcmStorageLock
	buf   []byte
	open  int

	cyp cipher.AEAD
}

// OpenEncryptedFile ...
func OpenEncryptedFile(path string,
	key []byte, readOnly bool) (storage.Storage, error) {
	ace, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cyp, err := cipher.NewGCM(ace)
	if err != nil {
		return nil, err
	}

	if fi, err := os.Stat(path); err == nil {
		if !fi.IsDir() {
			return nil, fmt.Errorf(
				"leveldb/storage: open %s: not a directory", path)
		}
	} else if os.IsNotExist(err) && !readOnly {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	flock, err := newFileLock(filepath.Join(path, "LOCK"), readOnly)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			flock.release()
		}
	}()

	fs := &aesgcmStorage{
		path:     path,
		readOnly: readOnly,
		flock:    flock,
		cyp:      cyp,
	}
	runtime.SetFinalizer(fs, (*aesgcmStorage).Close)
	return fs, nil
}

func (fs *aesgcmStorage) Log(_ string) {
}

func fdGenAD(fd storage.FileDesc) []byte {
	ret := make([]byte, additionalDataLen)
	ret[0] = byte(fd.Type)
	binary.LittleEndian.PutUint64(ret[1:], uint64(fd.Num))
	return ret
}

func (fs *aesgcmStorage) Open(fd storage.FileDesc) (storage.Reader, error) {
	fs.Log(fmt.Sprintf("opening %s", fd))
	if !storage.FileDescOk(fd) {
		return nil, storage.ErrInvalidFile
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.open < 0 {
		return nil, storage.ErrClosed
	}
	of, err := os.OpenFile(
		filepath.Join(fs.path, fsGenName(fd)), os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer of.Close()

	crypt, err := ioutil.ReadAll(of)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, fs.cyp.NonceSize())
	copy(nonce[0:fs.cyp.NonceSize()], crypt[0:fs.cyp.NonceSize()])
	plain, err := fs.cyp.Open(
		nil, nonce, crypt[fs.cyp.NonceSize():], fdGenAD(fd))
	if err != nil {
		return nil, err
	}
	fs.open++
	return newReader(plain, fd, fs), nil
}

func (fs *aesgcmStorage) Create(fd storage.FileDesc) (storage.Writer,
	error) {
	fs.Log(fmt.Sprintf("create %s", fd))

	if !storage.FileDescOk(fd) {
		return nil, storage.ErrInvalidFile
	}
	if fs.readOnly {
		return nil, errReadOnly
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.open < 0 {
		return nil, storage.ErrClosed
	}
	of, err := os.OpenFile(
		filepath.Join(fs.path, fsGenName(fd)),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	fs.open++
	return newWriter(of, fd, fs), nil
}

func (fs *aesgcmStorage) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.open < 0 {
		return storage.ErrClosed
	}
	runtime.SetFinalizer(fs, nil)

	if fs.open > 0 {
		fs.Log(fmt.Sprintf("close: warning, %d files still open", fs.open))
	}
	fs.open = -1
	return fs.flock.release()
}
