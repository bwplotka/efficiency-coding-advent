# efficiency-advent-2021

Repo for [Advent of Code 2021](https://adventofcode.com/2021) in Go with efficient solutions. Note that I focus on "maintainable" efficiency. Imagine those algorithms are needed in critical path, but normal humans have to still be able to maintain it! (: 

You will read more about tricks used here in [Efficient Go Book](https://www.oreilly.com/library/view/efficient-go/9781098105709/), but I left comments on various solutions through the challenges, enjoy! (: 

Inspired by [Advent of Go Profiling 2021](https://felixge.de/2021/12/01/advent-of-go-profiling-2021-day-1-1/.)

Follow others in their Advent of Code journey too!

* [@felixge](https://felixge.de/2021/12/01/advent-of-go-profiling-2021-day-1-1/)
* [@kabanek](https://twitter.com/kabanek/status/1466284532269821959)

## How to Use this Repo

Each day has two parts. Both are within the same `main.go` file within corresponding `dayXY` directory. You can find correctness tests and benchmarks within `main_test.go`

Benchmark typically looks as follows:

```go 
var Answer int

func BenchmarkFUNC_NAME(b *testing.B) {
	b.ReportAllocs()
	input := strings.TrimSpace(day2.ReadTestInput(b)) + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = FUNC_NAME(input)
	}
}
```

Notice the `Answer` var on global scope. This makes sure Go compiler will not optimize our function away.

You can run 5 tests, if you run command:

You can run this test by going into directory with this file and running:

`export var=v1 && go test -count 5 -run '^$' -bench . -memprofile=${var}.mem.pprof -cpuprofile=${var}.cpu.pprof > ${var}.txt`

This will give pprof profiles you can inspect by running:

`go tool pprof -edgefraction=0 -functions -http=:8081 ${var}.cpu.pprof`

or for memory:

`go tool pprof -edgefraction=0 -functions -http=:8081 ${var}.mem.pprof`

You can also find output of benchmark using `benchstat` you can install using `go install golang.org/x/perf/cmd/benchstat`. This tool will get statistically correct value for multiple runs:

`benchstat ${var}.txt`

You can also compare different results across each other e.g v1 and v5:

`benchstat v1.txt v5.txt`

That's it! Enjoy (:
