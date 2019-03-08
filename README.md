# `spanmap` 
[![Build Status](https://img.shields.io/travis/tidwall/spanmap.svg?style=flat-square)](https://travis-ci.org/tidwall/spanmap)
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tidwall/spanmap)
[![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](http://gocover.io/github.com/tidwall/spanmap)

A specialized collection type that uses `uint64` for the index key 
and `interface{}` for the data. It's ideal for entries that have indexes that 
are clustered together, such as transactions logs, time series, and other 
monotonically or contiguously ordered data.

You can think of `spanmap` like a cross between a map and an array. It can have 
any starting index, dynamically grow, and contain small gaps. 

*This collection should not be used for data where the indexes have
random properties or very large gaps.*

# Getting Started

### Installing

To start using spanmap, install Go and run `go get`:

```sh
$ go get -u github.com/tidwall/spanmap
```

## Example 

```go
var m spanmap.Map

m.Set(910003, "3")
m.Set(910001, "1")
m.Set(910004, "4")
m.Set(910002, "2")

for i := m.Min(); i <= m.Max(); i++ {
    println(i, m.Get(i).(string))
}
```

## All operations

```go
func (m *Map) Set(index uint64, value interface{}) interface{}
func (m *Map) Get(index uint64) interface{}
func (m *Map) Delete(index uint64) interface{}
func (m *Map) Len() int
func (m *Map) MinIndex() uint64
func (m *Map) MaxIndex() uint64
```

Outputs:

```
910001 1
910002 2
910003 3
910004 4
```

## Performance

Single threaded performance comparing this package to the 
stdlib `map[uint64]interface{}`, 
[google/btree](https://github.com/google/btree), and 
[tidwall/celltree](https://github.com/tidwall/celltree).

```
$ go test

-- spanmap --
set/sequential 1,000,000 ops in 65ms, 15,402,400/sec, 64 ns/op
get/sequential 1,000,000 ops in 3ms, 292,561,563/sec, 3 ns/op
set/random     1,000,000 ops in 130ms, 7,672,620/sec, 130 ns/op
get/random     1,000,000 ops in 29ms, 34,351,759/sec, 29 ns/op

-- stdlib map --
set/sequential 1,000,000 ops in 338ms, 2,957,562/sec, 338 ns/op
get/sequential 1,000,000 ops in 84ms, 11,881,738/sec, 84 ns/op
set/random     1,000,000 ops in 299ms, 3,339,927/sec, 299 ns/op
get/random     1,000,000 ops in 84ms, 11,944,527/sec, 83 ns/op

-- btree --
set/sequential 1,000,000 ops in 216ms, 4,623,031/sec, 216 ns/op
get/sequential 1,000,000 ops in 169ms, 5,924,168/sec, 168 ns/op
set/random     1,000,000 ops in 689ms, 1,451,352/sec, 689 ns/op
get/random     1,000,000 ops in 884ms, 1,131,258/sec, 883 ns/op

-- celltree --
set/sequential 1,000,000 ops in 80ms, 12,465,058/sec, 80 ns/op
get/sequential 1,000,000 ops in 104ms, 9,625,261/sec, 103 ns/op
set/random     1,000,000 ops in 254ms, 3,944,502/sec, 253 ns/op
get/random     1,000,000 ops in 318ms, 3,142,246/sec, 318 ns/op
```

*These benchmarks were run on a MacBook Pro 15" 2.9 GHz Intel Core i9 
using Go 1.12*

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

`spanmap` source code is available under the MIT [License](/LICENSE).
