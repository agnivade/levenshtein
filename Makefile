# FUNCS is the list of functions to run (all by default)
FUNCS?=./...

# COUNT is the number of times to run each benchmark (5 by default)
COUNT?=5

all: test install

install:
	go install

lint:
	gofmt -l -s -w . && go vet .

test:
	go test -race -v -coverprofile=coverage.out -covermode=atomic

# bench runs benchmarks on levenshtein function
bench:
	go test -bench="BenchmarkDistanceAgnivade" -benchmem -count=${COUNT} | tee bench_agnivade.out
	benchstat bench_agnivade.out

# benchAll compares performance of all levenshtein distance implementations
benchAll:
	go test -bench="BenchmarkDistanceAgnivade" -benchmem -count=${COUNT} | tee bench_agnivade.out
	go test -bench="BenchmarkDistanceArbovm" -benchmem -count=${COUNT} | tee bench_arbovm.out
	go test -bench="BenchmarkDistanceDgryski" -benchmem -count=${COUNT} | tee bench_dgryski.out
	benchstat -col=.name bench_agnivade.out bench_arbovm.out bench_dgryski.out | tee bench_all.out

# cover runs the tests and opens a web browser displaying annotated source code
cover: test
	@if [ $$? -eq 0 ]; then \
		go tool cover -html=coverage.out; \
	fi

# fuzz runs all fuzzing functions
fuzz:
	go test -fuzz="$(FUNCS)"

# perf compares performance using benchstat between the last commit and uncommitted code
# to install benchstat run 'go install golang.org/x/perf/cmd/benchstat@latest'
perf:
	go test -bench="BenchmarkDistanceAgnivade" -benchmem -count=${COUNT} | tee perf_after.out
	git stash -q --keep-index
	go test -bench="BenchmarkDistanceAgnivade" -benchmem -count=${COUNT} | tee perf_before.out
	git stash pop -q
	benchstat perf_before.out perf_after.out | tee perf_diff.out