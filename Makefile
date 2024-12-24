# FUNCS is the list of functions to run
FUNCS?=./...

# COUNT is the number of times to run each benchmark
COUNT?=5

all: test install

install:
	go install

lint:
	gofmt -l -s -w . && go vet .

test:
	go test -race -v -coverprofile=coverage.out -covermode=atomic

bench:
	go test -bench="$(FUNCS)" -benchmem -count=${COUNT}

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
	go test -bench="$(FUNCS)" -benchmem -count=${COUNT}  | tee perf_after.out
	git stash -q --keep-index
	go test -bench="$(FUNCS)" -benchmem -count=${COUNT} | tee perf_before.out
	git stash pop -q
	benchstat perf_before.out perf_after.out | tee perf_diff.out