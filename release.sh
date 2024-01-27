#! /usr/bin/env bash
set -uvx
set -e
cwd=`pwd`
version=0.1.6
git add .
git commit -m"Release v$version"
git tag -a v$version -mv$version
git push origin v$version
git push
git remote -v
gh auth login --with-token < ~/settings/github-all-tokne.txt
gh release create v$version --generate-notes --target main
