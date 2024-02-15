#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test archive command
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
arcdir="${tmpd}/archives"
mkdir -p "${arcdir}"

# create archive
tar_arc="${arcdir}/archive1.tar.gz"
tar -czf "${tar_arc}" "${cur}/../internal"

zip_arc="${arcdir}/archive2.zip"
zip -r "${zip_arc}" "${cur}/../internal"

# index
"${bin}" --debug index -a -C -c "${catalog}" "${arcdir}" arcdir
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

cnt=$(find "${cur}/../internal" | wc -l)
# +1 storage entry
# +2 each of the archive file header
total="$(("${cnt}" + "${cnt}" + 1 + 2))"

echo ">>> test archive ls <<<"
"${bin}" --debug ls -r -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${total}" ] && echo "expecting ${total} line (got ${cnt})" && exit 1
grep '^storage arcdir 0B' "${out}" && (echo "empty storage" && exit 1)
grep '^archive1.tar.gz' "${out}" || (echo "empty storage" && exit 1)
grep '^archive2.zip' "${out}" || (echo "empty storage" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
