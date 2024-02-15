#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6

cur=$(cd "$(dirname "${0}")" && pwd)

# integration tests
find "${cur}" -iname 'test-*.sh' | while read -r line; do
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

echo "all tests OK!"