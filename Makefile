.PHONY: default get codetest test fmt lint vet

default: codetest

get:
	go get -v ./...
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.20.0

codetest: fmt lint vet test

test:
	go test -v -cover

fmt:
	go fmt ./...

lint:
	golangci-lint run --fix

vet:
	go vet -all .
