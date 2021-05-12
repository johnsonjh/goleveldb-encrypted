# JLevelDB Encrypted Storage

[![GoReportCard](https://goreportcard.com/badge/github.com/johnsonjh/jleveldb-encrypted)](https://goreportcard.com/report/github.com/johnsonjh/jleveldb-encrypted)

JLevelDB Encrypted Storage provides AES encryption support for
[JLevelDB](https://github.com/johnsonjh/jleveldb).

## Installation

- `go get -a -v github.com/johnsonjh/jleveldb-encrypted`

## Usage

A wrapper equivalent to `OpenFile` is provided. 

Replace calls to `OpenFile` with `OpenAESEncryptedFile`.
```
db, err = OpenAESEncryptedFile( dir, key, nil )
defer db.Close()
db.Put( []byte("hello"), []byte("value") )
```

## Performance

```
Linux/amd64, Go 1.17-cb88c5b6be, Intel Core i7-8700

Benchmark_Normal_100keys_8bytes-12       	     859	   1433694 ns/op
Benchmark_AES128_100keys_8bytes-12       	     759	   1470218 ns/op
Benchmark_AES256_100keys_8bytes-12       	     808	   1480387 ns/op

Benchmark_Normal_10000keys_8bytes-12     	      55	  21508465 ns/op
Benchmark_AES128_10000keys_8bytes-12     	      64	  17735344 ns/op
Benchmark_AES256_10000keys_8bytes-12     	     418	   2834664 ns/op

Benchmark_Normal_100keys_32bytes-12      	     837	   1463877 ns/op
Benchmark_AES128_100keys_32bytes-12      	     795	   1515165 ns/op
Benchmark_AES256_100keys_32bytes-12      	     789	   1506116 ns/op

Benchmark_Normal_10000keys_32bytes-12    	      50	  23339237 ns/op
Benchmark_AES128_10000keys_32bytes-12    	      56	  20262302 ns/op
Benchmark_AES256_10000keys_32bytes-12    	      58	  20321537 ns/op
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


