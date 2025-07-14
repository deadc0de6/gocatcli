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

# add some hidden and non hidden files
mkdir "${src}/.hidden"
echo "hidden" > "${src}/.hidden/hiddenfile"
mkdir "${src}/nothidden"
echo "not-hidden" > "${src}/nothidden/subfile"
echo "hidden" > "${src}/.file-hidden"
echo "not-hidden" > "${src}/notfile-hidden"

# add some with extensions
mkdir -p "${src}/with.ext"
echo "withext" > "${src}/with.ext/inside"
echo "theext" > "${src}/file.ext"

# index
echo ">>> indexing <<<"
"${bin}" index -a -C -c "${catalog}" --debug --ignore='\.+' --ignore='\.ext' "${tmpd}/to-index" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
echo ">>> test index ignore ls <<<"
"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#expected=$(find "${cur}/../" -not -path '*/.git*' | grep -v '^.$' | wc -l)
cat_file "${out}"

expected=4
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

grep 'notfile-hidden' "${out}" || exit 1
grep 'nothidden' "${out}" || exit 1
grep 'subfile' "${out}" || exit 1
grep '\.hidden' "${out}" && exit 1
grep 'hiddenfile' "${out}" && exit 1
grep '\.file-hidden' "${out}" && exit 1

echo "test $(basename "${0}") OK!"
exit 0
