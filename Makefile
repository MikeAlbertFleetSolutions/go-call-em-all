.PHONY: default get codetest fmt lint vet

default: fmt codetest

get:
ifneq ("$(CI)", "true")
	go get -u ./...
	go mod tidy
endif
	go mod download
	go mod verify

codetest: lint vet

fmt:
	go fmt ./...

lint:
ifeq ("$(CI)", "true")
	$(shell go env GOPATH)/bin/golangci-lint run --verbose --timeout 3m
else
	$(shell go env GOPATH)/bin/golangci-lint run --fix
endif

vet:
	go vet -all ./...
