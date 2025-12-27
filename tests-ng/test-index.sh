#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test index command
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
"${bin}" index -a -C -c "${catalog}" --debug --ignore='**/.git*/**' "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
echo ">>> test index ls <<<"
"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#expected=$(find "${cur}/../" -not -path '*/.git*' | grep -v '^.$' | wc -l)
cat_file "${out}"
expected=$("${cur}/plist.py" "${cur}/../" --ignore '*/.git*')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

catalog="${tmpd}/catalog2"

# index
echo ">>> test index <<<"
"${bin}" index -a -C -c "${catalog}" "${cur}/../internal" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
echo ">>> test index ls (2) <<<"
"${bin}" --debug -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#expected=$(find "${cur}/../internal" -not -path '*/.git*' | grep -v '^.$' | wc -l)
expected=$("${cur}/plist.py" "${cur}/../internal" --ignore '*/\.git*')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

echo "test $(basename "${0}") OK!"
exit 0
