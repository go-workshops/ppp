Bench time: 10s

```text
goos: darwin
goarch: arm64
pkg: github.com/go-workshops/ppp/benchmarks/hot_path
BenchmarkWork1-10          82309            168573 ns/op            5338 B/op        123 allocs/op        92410 DB records
BenchmarkWork2-10          71518            171544 ns/op            5606 B/op        128 allocs/op        45863 DB records
PASS
ok      github.com/go-workshops/ppp/benchmarks/hot_path 46.007s
```

```text
goos: darwin
goarch: arm64
pkg: github.com/go-workshops/ppp/benchmarks/hot_path
BenchmarkWork1-10          83258            135293 ns/op            5330 B/op        123 allocs/op        93359 DB records
BenchmarkWork2-10          86157            147063 ns/op            5587 B/op        128 allocs/op        53183 DB records
PASS
ok      github.com/go-workshops/ppp/benchmarks/hot_path 41.573s
```

```text
goos: darwin
goarch: arm64
pkg: github.com/go-workshops/ppp/benchmarks/hot_path
BenchmarkWork1-10          82740            139374 ns/op            5335 B/op        123 allocs/op        92841 DB records
BenchmarkWork2-10          83820            145393 ns/op            5589 B/op        128 allocs/op        52014 DB records
PASS
ok      github.com/go-workshops/ppp/benchmarks/hot_path 42.737s
```
