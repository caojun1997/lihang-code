.PHONY: build run test clean deps

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server/main.go

test:
	go test -v ./...

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/
	rm -rf uploads/*
	rm -rf results/*

lint:
	golangci-lint run

fmt:
	go fmt ./...
	gofmt -s -w .

docker-build:
	docker build -t pdf-parser:latest -f Dockerfile .

docker-run:
	docker run -p 8080:8080 pdf-parser:latest
