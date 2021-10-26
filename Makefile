.PHONY: all test run clean

all:
	go build -o pwdhash main.go

test:
	@go test -count 1 -cover ./...

run: all
	./pwdhash

clean:
	@rm -f pwdhash
