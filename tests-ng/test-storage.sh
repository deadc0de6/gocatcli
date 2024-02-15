#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test storage command
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
echo ">>> test index <<<"
"${bin}" index -a -C -c "${catalog}" "${cur}/../internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

"${bin}" index -a -C -c "${catalog}" "${cur}/../tests-ng" testsng
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

echo ">>> test list <<<"
"${bin}" --debug -c "${catalog}" storage list | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="2"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1

echo ">>> test meta <<<"
"${bin}" --debug -c "${catalog}" storage meta internal XYZ | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="2"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep "meta:XYZ" "${out}" || (echo "meta not saved" && exit 1)

echo ">>> test tag add 1 <<<"
"${bin}" --debug -c "${catalog}" storage tag testsng tag1 | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="2"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && (echo "expecting ${expected} lines (${cnt})" && exit 1)
grep "tags:tag1" "${out}" || (echo "tag1 not saved" && exit 1)

echo ">>> test tag add 2 <<<"
"${bin}" --debug -c "${catalog}" storage tag testsng tag2 | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="2"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && (echo "expecting ${expected} lines (${cnt})" && exit 1)
grep "tags:tag1,tag2" "${out}" || (echo "tag2 not saved" && cat "${out}" && exit 1)

echo ">>> test untag <<<"
"${bin}" --debug -c "${catalog}" storage untag testsng tag1 | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="2"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && (echo "expecting ${expected} lines (${cnt})" && exit 1)
grep "tags:tag2" "${out}" || (echo "tag1 not removed" && exit 1)

echo ">>> test rm storage <<<"
"${bin}" --debug -c "${catalog}" storage rm -f testsng | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="1"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep "^storage internal" "${out}" || (echo "testsng not removed" && exit 1)

echo ">>> test list <<<"
"${bin}" --debug -c "${catalog}" storage list | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected="1"
cat "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && (echo "expecting ${expected} lines (${cnt})" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
