`jsonx` is parser for arbitrary json in Golang, the `x` stands for arbitrary. The big goal of this package is to be more performant than the standard library, for the same task. Currently `jsonx` is more than twice as fast as `json.Unmarshal`.

Benchmarks
```
goos: windows
goarch: amd64
pkg: github.com/hk-32/jsonx/bench
cpu: 13th Gen Intel(R) Core(TM) i5-13400F
BenchmarkCustom-16      	     562	   2092281 ns/op	 2247421 B/op	   22024 allocs/op
BenchmarkStandard-16    	     259	   4606611 ns/op	 2417647 B/op	   55557 allocs/op
```