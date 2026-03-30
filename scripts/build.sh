#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o aquestalk-server cmd/aquestalk-server/main.go
