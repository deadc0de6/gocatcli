#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test size command
#

## start-test-cookie
set -eu -o errtrace -o pipefail
cur=$(cd "$(dirname "${0}")" && pwd)
bin="${cur}/../bin/gocatcli"
[ ! -e "${bin}" ] && echo "\"${bin}\" not found" && exit 1
# shellcheck disable=SC1091
source "${cur}"/helpers
## end-test-cookie

######################################
## the test

tmpd=$(mktemp -d --suffix='-dotdrop-tests' || mktemp -d)
clear_on_exit "${tmpd}"

catalog="${tmpd}/catalog"
out="${tmpd}/output.txt"

# index
"${bin}" --debug index -a -C -c "${catalog}" "${cur}/../internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

echo ">>> test file size <<<"
"${bin}" --debug ls -r -l -S -c "${catalog}" "${cur}/../internal/walker/walker.go" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#expected=$(du --block=1 --apparent-size "${cur}/../internal/walker/walker.go" | awk '{print $1}')
# shellcheck disable=SC2012
expected=$(ls -la "${cur}/../internal/walker/walker.go" | awk '{print $5}')
cat_file "${out}"
size=$(grep 'walker.go' "${out}" | awk '{print $4}')
echo "size:${size} VS exp:${expected}"
[ "${size}" != "${expected}" ] && echo "expecting ${expected} (got ${size})" && exit 1

echo ">>> test directory size <<<"
"${bin}" --debug ls -l -r -S -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#expected=$(du -c --block=1 --apparent-size "${cur}/../internal/catcli" | tail -1 | awk '{print $1}')
"${cur}/pdu.py" "${cur}/../internal/catcli" | tail -1
expected=$("${cur}/pdu.py" "${cur}/../internal/catcli" | tail -1 | awk '{print $1}')
#cat_file "${out}"
size=$(grep '^catcli' "${out}" | awk '{print $4}')
echo "size:${size} VS exp:${expected}"
[ "${size}" != "${expected}" ] && echo "expecting ${expected} (got ${size})" && exit 1

echo ">>> test storage size <<<"
"${bin}" --debug ls -l -S -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
#expected=$(du -c --block=1 --apparent-size "${cur}/../internal" | tail -1 | awk '{print $1}')
expected=$("${cur}/pdu.py" "${cur}/../internal" | tail -1 | awk '{print $1}')
size=$(grep '^storage' "${out}" | awk '{print $3}')
echo "size:${size} VS exp:${expected}"
[ "${size}" != "${expected}" ] && echo "expecting ${expected} (got ${size})" && exit 1

echo "test $(basename "${0}") OK!"
exit 0
