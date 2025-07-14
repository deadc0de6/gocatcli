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

## create a fake dir
src="${tmpd}/to-index"
mkdir -p "${src}"
mkdir "${src}/.hidden"
echo "hidden" > "${src}/.hidden/hiddenfile"
mkdir "${src}/nothidden"
echo "not-hidden" > "${src}/nothidden/subfile"
echo "hidden" > "${src}/.file-hidden"
echo "not-hidden" > "${src}/notfile-hidden"

# index
echo ">>> indexing <<<"
"${bin}" index -a -C -c "${catalog}" --ignore='.*' "${tmpd}/to-index" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
echo ">>> test ls <<<"
"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#expected=$(find "${cur}/../" -not -path '*/.git*' | grep -v '^.$' | wc -l)
cat_file "${out}"

expected=4
cnt=$(cat "${out}" | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

cat "${out}" | grep 'notfile-hidden' || exit 1
cat "${out}" | grep 'nothidden' || exit 1
cat "${out}" | grep 'subfile' || exit 1

cat "${out}" | grep '\.hidden' && exit 1
cat "${out}" | grep 'hiddenfile' && exit 1
cat "${out}" | grep '\.file-hidden' && exit 1

echo "test $(basename "${0}") OK!"
exit 0
