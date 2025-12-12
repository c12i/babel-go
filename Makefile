.PHONY: fmt vet lint test

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test -v ./...

check: fmt vet lint test
