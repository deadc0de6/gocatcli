#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test find command
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

# find no arg
echo ">>> test find no arg <<<"
"${bin}" find -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../internal" ! -path '*/.git/*' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep '^storage internal.*' "${out}" || (echo "bad content" && exit 1)

echo ">>> test find no star <<<"
"${bin}" --debug find -c "${catalog}" nav | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../internal" ! -path '*/.git/*' -name 'nav*' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep '^internal/commands/nav.go.*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/navigator.*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/navigator/navigator.go.*' "${out}" || (echo "bad content" && exit 1)

echo ">>> test find with star <<<"
"${bin}" --debug find -c "${catalog}" '*nav*' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../internal" ! -path '*/.git/*' -name '*nav*' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep '^internal/commands/nav.go.*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/navigator.*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/navigator/navigator.go.*' "${out}" || (echo "bad content" && exit 1)

echo ">>> test find with only leading star <<<"
"${bin}" --debug find -c "${catalog}" '*gator' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../internal" ! -path '*/.git/*' -name '*gator' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1

echo ">>> test find with only trailing star <<<"
"${bin}" --debug find -c "${catalog}" 'find*' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../internal" ! -path '*/.git/*' -name 'find*' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1

echo ">>> test find with multiple args <<<"
"${bin}" --debug find -c "${catalog}" storage tree | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../internal" | grep -c 'storage\|tree')
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1
grep '^internal/commands/storage.go.*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/node/storage.go .*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/tree .*' "${out}" || (echo "bad content" && exit 1)
grep '^internal/tree/tree.go' "${out}" || (echo "bad content" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
