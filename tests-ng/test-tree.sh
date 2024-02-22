#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test tree command
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
"${bin}" index -a -C -c "${catalog}" --ignore=".git" "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# tree
echo ">>> test tree no arg <<<"
"${bin}" --debug tree -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
#expected=$(find "${cur}/../" -not -path '*/.git*' | tail -n +2 | wc -l)
cat_file "${out}"
expected=$("${cur}/plist.py" "${cur}/../" --ignore '*/.git*')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

# tree with arg
echo ">>> test tree with arg <<<"
"${bin}" --debug tree -a -c "${catalog}" internal | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
expected=$("${cur}/plist.py" "${cur}/../internal" --ignore '*/.git*')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

echo "test $(basename "${0}") OK!"
exit 0
