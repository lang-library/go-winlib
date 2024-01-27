#! /usr/bin/env bash
set -uvx
set -e
rm -rf go.mod go.sum
go mod init github.com/lang-library/go-winlib
go get golang.org/x/sys/windows/mkwinsyscall
go mod tidy
cat go.mod
