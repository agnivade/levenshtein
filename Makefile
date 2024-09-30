all: test install

install:
	go install

lint:
	gofmt -l -s -w . && go vet .

test:
	go test -race -v -coverprofile=coverage.out -covermode=atomic

bench:
	go test -run=XXX -bench=. -benchmem -count=5

# cover runs the tests and opens a web browser displaying annotated source code
cover: test
	@if [ $$? -eq 0 ]; then \
		go tool cover -html=coverage.out; \
	fi

# fuzz run all fuzzing functions
fuzz:
	go test -fuzz=.