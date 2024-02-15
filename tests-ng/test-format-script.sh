#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test format script command
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
"${bin}" index -C -c "${catalog}" "${cur}/../internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

#"${bin}" tree -c "${catalog}"
#"${bin}" ls -r -c "${catalog}"

"${bin}" ls -r -S -a --format=script -c "${catalog}" internal/tree | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"

# shellcheck disable=SC2016
exp='op=file; source=/media/mnt; ${op} "${source}/tree" "${source}/tree/tree.go"'
grep "${exp}" "${out}" || (echo "bad output" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
