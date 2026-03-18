default: fmt lint build generate_docs

dev:
	goreleaser build --id $(shell go env GOOS) --single-target --snapshot --clean

snapshot:
	goreleaser release --snapshot --clean

build: dev

lint:
	golangci-lint run

generate_docs:
	go tool tfplugindocs

fmt:
	gofmt -s -w -e .

testacc:
	go test -v -cover -timeout=120s -parallel=10 ./...

.PHONY: dev snapshot build lint generate_docs fmt testacc
