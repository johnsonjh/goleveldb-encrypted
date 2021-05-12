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
 */

package aesgcm

import (
	"bytes"

	"github.com/johnsonjh/jleveldb/leveldb/storage"
)

type aesgcmReader struct {
	*bytes.Reader
	fs     *aesgcmStorage
	fd     storage.FileDesc
	closed bool
}

func newReader(b []byte, fd storage.FileDesc,
	fs *aesgcmStorage) *aesgcmReader {
	return &aesgcmReader{
		Reader: bytes.NewReader(b),
		fs:     fs,
		fd:     fd,
		closed: false,
	}
}

func (r aesgcmReader) Close() error {
	r.fs.mu.Lock()
	defer r.fs.mu.Unlock()
	if r.closed {
		return storage.ErrClosed
	}
	r.closed = true
	r.fs.open--
	return nil
}
