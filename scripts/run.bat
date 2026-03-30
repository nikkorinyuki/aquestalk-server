@echo off

set GOOS=linux
set GOARCH=amd64
go run cmd/aquestalk-server/main.go
