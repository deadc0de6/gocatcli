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

# create archives
tar_arc="${arcdir}/archive1.tar.gz"
tar -C "${cur}/.." -czf "${tar_arc}" "internal"

zip_arc="${arcdir}/archive2.zip"
(cd "${cur}/../" && zip -r "${zip_arc}" "internal")

echo "test" > "${tmpd}/gzipped-original"
gzip -c "${tmpd}/gzipped-original" > "${arcdir}/gzipped.gz"

unzip -l "${zip_arc}"

# index
"${bin}" --debug index -a -C -c "${catalog}" "${arcdir}" arcdir
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

cnt=$(find "${cur}/../internal" | wc -l)
# +1 storage entry
# +2 each of the archive file header
# +1 for gzipped
total="$(("${cnt}" + "${cnt}" + 1 + 2 + 1))"

echo ">>> test archive ls <<<"
"${bin}" --debug ls -r -a -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
cat_file "${out}"
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${total}" ] && echo "expecting ${total} line (got ${cnt})" && exit 1
grep '^storage arcdir 0B' "${out}" && (echo "empty storage" && exit 1)
grep '^archive1.tar.gz' "${out}" || (echo "no archive1" && exit 1)
grep '^archive2.zip' "${out}" || (echo "archive2" && exit 1)
grep '^gzipped.gz' "${out}" || (echo "no gzipped" && exit 1)

echo ">>> test archive create <<<"
dst="${tmpd}/created"
"${bin}" --debug -c "${catalog}" create --archive "${dst}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"

tree "${dst}/arcdir"
#cat "${catalog}"
[ ! -d "${dst}/arcdir" ] && echo "top not created" && exit 1
[ ! -d "${dst}/arcdir/archive1.tar.gz" ] && (echo "no create archive1" && exit 1)
[ ! -d "${dst}/arcdir/archive2.zip" ] && (echo "no create archive2" && exit 1)
[ ! -f "${dst}/arcdir/gzipped.gz" ] && (echo "no create gzipped" && exit 1)
[ ! -f "${dst}/arcdir/archive1.tar.gz/internal/tree/tree.go" ] && (echo "no create tree.go" && exit 1)
[ ! -f "${dst}/arcdir/archive1.tar.gz/internal/catcli/convertor.go" ] && (echo "no create convertor.go" && exit 1)
[ ! -f "${dst}/arcdir/archive1.tar.gz/internal/walker/archives/archive.go" ] && (echo "no create archive.go" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
