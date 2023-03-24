run: build exe

build:
	@go build -o ./bin/out
exe:
	@./bin/out
