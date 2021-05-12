# JLevelDB Encrypted Storage

[![GoReportCard](https://goreportcard.com/badge/github.com/johnsonjh/jleveldb-encrypted)](https://goreportcard.com/report/github.com/johnsonjh/jleveldb-encrypted)

JLevelDB Encrypted Storage provides a strongly encrypted storage
(data at rest) for [JLevelDB](https://github.com/johnsonjh/jleveldb).

## Installation

- `go get -a -v github.com/johnsonjh/jleveldb-encrypted`

## Usage

The storage engine can be manually instantiated (see the `aesgcm`
package for the raw storage interface), but for most use cases a
wrapper equivalent to `OpenFile` is provided - replace calls to
`OpenFile` with calls to `OpenAESEncryptedFile`.
```
db, err = OpenAESEncryptedFile( dir, key, nil )
defer db.Close()
db.Put( []byte("hello"), []byte("value") )
```

## Performance

```
Linux/amd64 - Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz

Benchmark_Normal_100keys_8bytes-12       	     753	   1627542 ns/op
Benchmark_AES128_100keys_8bytes-12       	     723	   1666636 ns/op
Benchmark_AES256_100keys_8bytes-12       	     715	   1664757 ns/op
Benchmark_Normal_10000keys_8bytes-12     	      52	  22408747 ns/op
Benchmark_AES128_10000keys_8bytes-12     	      66	  17912655 ns/op
Benchmark_AES256_10000keys_8bytes-12     	     405	   3045821 ns/op
Benchmark_Normal_100keys_32bytes-12      	     760	   1656028 ns/op
Benchmark_AES128_100keys_32bytes-12      	     630	   1746463 ns/op
Benchmark_AES256_100keys_32bytes-12      	     679	   1748278 ns/op
Benchmark_Normal_10000keys_32bytes-12    	      48	  23738455 ns/op
Benchmark_AES128_10000keys_32bytes-12    	      56	  21215759 ns/op
Benchmark_AES256_10000keys_32bytes-12    	      57	  20984932 ns/op
```

## License

This package contains code from the GoLevelDB and GoLevelDB Encrypted Storage
projects.

JLevelDB Encrypted Storage is licensed under the Apache License, Version 2.0
(the "License"); you may not use this package except in compliance with the
License. You may obtain a copy of the License at
[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0).
Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

In addition, this entire repository may also be used under the 3-Clause BSD
license of [GoLevelDB Project](https://github.com/syndtr/goleveldb/blob/master/LICENSE).

## Original Authors

Derived from [GoLevelDB Encrypted Storage](https://github.com/tenta-browser/goleveldb-encrypted) by [Tenta](https://tenta.com).


