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

package aesgcm_test

import (
	"math/rand"
	"os"
	"testing"

	"github.com/johnsonjh/jleveldb-encrypted/aesgcm"
	u "github.com/johnsonjh/leaktestfe"
)

func TestOpenEncryptedFile_Keys(t *testing.T) {
	defer u.Leakplug(t)
	for i := 0; i < 65; i++ {
		temp := tempDir(t)

		key := make([]byte, i)
		rand.Read(key)
		db, err := aesgcm.OpenEncryptedFile(temp, key, false)

		if i == 16 || i == 24 || i == 32 {
			if err != nil || db == nil {
				t.Logf("Should succeed when key length is %d", i)
				t.Fail()
			}
		} else {
			if err == nil || db != nil {
				t.Logf("Should fail when key length is %d", i)
				t.Fail()
			}
		}
		os.RemoveAll(temp)
	}
}
