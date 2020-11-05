#!/usr/bin/make -f

build: fmt
	go build -o bin/ ./...

fmt:
	go fmt ./...
