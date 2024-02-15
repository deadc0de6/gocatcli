#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test update command
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

cp -r "${cur}/../internal" "${tmpd}"

# index
echo ">>> test index <<<"
"${bin}" index -C -c "${catalog}" "${tmpd}/internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# refactor
rm -rf "${tmpd}/internal/tree"
rm -rf "${tmpd}/internal/utils"
mkdir -p "${tmpd}/internal/new"
echo "this is a test" > "${tmpd}/internal/new/newfile"

echo ">>> test re-index <<<"
"${bin}" index -a -f -C -c "${catalog}" "${tmpd}/internal" internal

#"${bin}" tree -c "${catalog}"
#"${bin}" ls -r -c "${catalog}"

"${bin}" ls -r -S -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
grep '^tree ' "${out}" && (echo "tree found" && exit 1)
grep '^utils ' "${out}" && (echo "utils found" && exit 1)
grep '^new' "${out}" || (echo "new not found" && exit 1)
grep 'newfile ' "${out}" || (echo "new not found" && exit 1)

expected=$(du -c --block=1 --apparent-size "${tmpd}/internal" | tail -1 | awk '{print $1}')
size=$(grep '^storage' "${out}" | awk '{print $3}')
echo "size:${size} VS exp:${expected}"
[ "${size}" != "${expected}" ] && echo "expecting ${expected} (got ${size})" && exit 1

echo ">>> test re-index with ignore <<<"
"${bin}" index -a -f -C -c "${catalog}" --ignore='*.go' "${tmpd}/internal" internal
"${bin}" ls -r -S -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
grep '^.*.go$' "${out}" && (echo ".go files found" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
