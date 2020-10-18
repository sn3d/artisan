#!/bin/sh
go mod download
go test -v ./...
go build -o bin/artisan