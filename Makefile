.PHONY: build test e2e dev clean lint benchmark bench profile race

BINARY = agentbridge

build:
	go build -o $(BINARY) .

test:
	go test ./... -v

e2e: build
	./$(BINARY) --help
	./$(BINARY) run "find top AI infra startups" --mode deterministic
	./$(BINARY) explain --task "find top AI infra startups"
	./$(BINARY) simulate-failure --seed 42

dev: build
	./$(BINARY) --help
	./$(BINARY) run "dev smoke task" --mode normal

lint:
	go vet ./...
	go fmt ./...

benchmark:
	go test -bench=. -benchmem -count=3 ./...

bench: benchmark

profile:
	go test -bench=BenchmarkEngineRunTaskDeterministic -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof -count=1 ./internal/workflow/...

race:
	go test -race ./...

clean:
	rm -f $(BINARY)
