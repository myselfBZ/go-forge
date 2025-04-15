do:
	@go build -o test/bin ./cmd

run:
	@go build -o bin/main ./cmd/
	@./bin/main

build:
	@go build -o bin/go-forge ./cmd/

clean:
	@rm -rf test
