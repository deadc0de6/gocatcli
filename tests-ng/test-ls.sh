#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test ls command
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

# ls no arg
echo ">>> test ls no arg <<<"
"${bin}" --debug ls -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "1" ] && echo "expecting single line (${cnt})" && exit 1
grep '^storage internal.*' "${out}" || (echo "bad content" && exit 1)

# ls no arg but rec
echo ">>> test ls no arg but rec <<<"
"${bin}" --debug ls -r -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
# shellcheck disable=SC2126
expected=$(find "${cur}/../internal" -not -path '*/.git*' | grep -v '^.$' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep '^storage internal.*' "${out}" || (echo "bad content" && exit 1)

# ls storage
echo ">>> test ls storage <<<"
"${bin}" --debug ls -l -a -c "${catalog}" internal | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "11" ] && echo "expecting 11 lines got ${cnt}" && exit 1
#grep '^storage internal.*' "${out}" || (echo "bad content 1" && exit 1)
grep 'fuser *d.*' "${out}" || (echo "bad content 2" && exit 1)
grep 'walker *d.*' "${out}" || (echo "bad content 3" && exit 1)

# ls recursive
echo ">>> test ls recursive <<<"
"${bin}" --debug ls -l -r -a -c "${catalog}" internal | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
# shellcheck disable=SC2126
expected=$(find "${cur}/../internal" -not -path '*/.git*' | grep -v '^.$' | tail -n+2 | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1
grep 'walker.go *-.*' "${out}" || (echo "bad content 1" && exit 1)
grep 'walker *d.*' "${out}" || (echo "bad content 2" && exit 1)

# ls pattern
echo ">>> test ls pattern <<<"
# expect to list
# - internal/navigator
# - internal/node
"${bin}" --debug ls -r -a -l -c "${catalog}" 'inter*/n*' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
expected=$(find "${cur}/../internal/navigator/"* "${cur}/../internal/node/"* | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1
grep 'node.go *-.*' "${out}" || (echo "bad content 2" && exit 1)
grep 'navigator.go *-.*' "${out}" || (echo "bad content 4" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
