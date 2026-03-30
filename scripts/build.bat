@echo off

set GOOS=linux
set GOARCH=amd64
go build -o aquestalk-server cmd/aquestalk-server/main.go
