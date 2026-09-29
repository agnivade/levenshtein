all: test install

install:
	go install

lint:
	gofmt -l -s -w . && go vet .

test:
	go test -race -v -coverprofile=coverage.txt -covermode=atomic

bench:
	cd benchmark && go test -run=XXX -bench=. -benchmem -count=5
