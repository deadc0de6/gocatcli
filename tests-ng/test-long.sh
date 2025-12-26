#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test long for ls and tree command
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

echo ">>> test ls <<<"
"${bin}" --debug ls -l -r -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
expected=$(find "${cur}/../internal" -not -path '*/.git*' | grep -v '^.$' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (got ${cnt})" && exit 1
grep 'indexed:' "${out}" || (echo "indexed not shown" && exit 1)
grep 'checksum:' "${out}" || (echo "checksum not shown" && exit 1)

echo ">>> test tree <<<"
"${bin}" --debug tree -l -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
expected=$(find "${cur}/../internal" -not -path '*/.git*' | grep -v '^.$' | wc -l)
# +1 empty line
expected=$((expected + 1))
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (got ${cnt})" && exit 1
grep 'indexed:' "${out}" || (echo "indexed not shown" && exit 1)
grep 'checksum:' "${out}" || (echo "checksum not shown" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
