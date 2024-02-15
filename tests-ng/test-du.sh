#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test du command
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
"${bin}" index -a -C -c "${catalog}" "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

echo ">>> test du raw <<<"
"${bin}" --debug du -S -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"

# shellcheck disable=SC2126
expected=$(find "${cur}/../" -type d | grep -v '^.$' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1

# bin size
echo ">>> test du bin size raw <<<"
echo "--- 1 ---"
du -c --block=1 --apparent-size "${cur}/../cmd/gocatcli"
echo "--- 2 ---"
du -c --block=1 "${cur}/../cmd/gocatcli"
echo "--- 3 ---"
du -c --apparent-size "${cur}/../cmd/gocatcli"
echo "--- 4 ---"
du -c "${cur}/../cmd/gocatcli"
echo "--- 5 ---"
ls -lah "${cur}/../cmd/gocatcli"
expected=$(du -c --block=1 --apparent-size "${cur}/../cmd/gocatcli" | tail -1 | awk '{print $1}')
size=$(grep '^.* *gocatcli/cmd/gocatcli$' "${out}" | awk '{print $1}')
echo "size:${size} VS exp:${expected}"
[ "${expected}" != "${size}" ] && (echo "bad bin size" && exit 1)

# total size
echo ">>> test du total size raw <<<"
expected=$(du -c --block=1 --apparent-size "${cur}/../" | tail -1 | awk '{print $1}')
#cat_file "${out}"
size=$(tail -1 "${out}" | awk '{print $1}')
echo "size:${size} VS exp:${expected}"
[ "${expected}" != "${size}" ] && (echo "bad total raw size" && exit 1)

echo ">>> test du human size <<<"
"${bin}" --debug du -c "${catalog}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
expected=$(find "${cur}/../" -type d | grep -v '^.$' | wc -l)
cnt=$(wc -l "${out}" | awk '{print $1}')
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines (${cnt})" && exit 1

# total size
echo ">>> test du total size human <<<"
expected=$(du -c --block=1 --apparent-size "${cur}/../" | tail -1 | awk '{print $1}' | sed 's/M//g')
# for some reason "du -h" uses 1000 with above options instead of 1024
expected=$(awk 'BEGIN {printf "%.0f",'"${expected}"'/1024/1024}')
cat_file "${out}"
size=$(tail -1 "${out}" | awk '{print $1}' | sed 's/MiB//g')
echo "size:${size} VS exp:${expected}"
[ "${expected}" != "${size}" ] && (echo "bad total human size" && exit 1)

echo "test $(basename "${0}") OK!"
exit 0
