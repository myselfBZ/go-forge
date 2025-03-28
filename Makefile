do:
	@go build -o test/bin .

run:
	@go build -o bin/main ./cmd/
	@./bin/main

clean:
	@rm -rf test
