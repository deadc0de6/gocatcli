#!/bin/bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6

set -eu -o errtrace -o pipefail
cur=$(cd "$(dirname "${0}")" && pwd)

# deps
echo "install deps..."
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/mgechev/revive@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# linting
echo "go fmt..."
go fmt ./...
#echo "golint..."
#golint -set_exit_status ./...
echo "golangci-lint..."
golangci-lint run

echo "revive..."
revive -set_exit_status -config ./revive.toml ./...

echo "staticcheck..."
staticcheck ./...

echo "go vet..."
go vet ./...

# unittests
#go test -v ./...

# shell scripts linting
find . -iname '*.sh' | while read -r script; do
  echo "-> checking \"${script}\""
  shellcheck "${script}"
done

# compilation
echo "compiling..."
make clean
make

# integration tests
find "${cur}/tests-ng" -iname 'test-*.sh' | while read -r line; do
  echo "running test script: ${line}"
  set +e
  "${line}"
  ret="$?"
  set -e
  if [ "${ret}" != "0" ]; then
    echo "test \"${line}\" failed!"
    exit 1
  fi
done

echo "done - all tests OK!"
