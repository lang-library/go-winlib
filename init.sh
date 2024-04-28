#! /usr/bin/env bash
set -uvx
set -e
get_latest_release () {
    gh release list --repo $1 | head -1 | sed 's/|/ /' | awk '{print $1, $8}'
}
gh auth login --with-token < ~/settings/github-all-tokne.txt
rm -rf go.mod go.sum
go mod init github.com/lang-library/go-winlib
go get golang.org/x/sys/windows/mkwinsyscall
go get github.com/lang-library/go-global@`get_latest_release lang-library/go-global`
go mod tidy
cat go.mod
