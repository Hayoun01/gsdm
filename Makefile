build:
	@go build -o ./bin/sdm

run: build
	@./bin/sdm $(ARGS)