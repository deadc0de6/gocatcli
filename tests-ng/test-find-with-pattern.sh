#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test find command with pattern
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

tmpd=$(mktemp -d --suffix='-gocatcli-tests' || mktemp -d)
clear_on_exit "${tmpd}"

catalog="${tmpd}/catalog"
out="${tmpd}/output.txt"

# index
"${bin}" index -a -C -c "${catalog}" "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ===================================================
title ">>> test find go files <<<"
"${bin}" find -c "${catalog}" '**/*.go' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../" -mindepth 1 -name '*.go' | wc -l | awk '{print $1}')
cnt=$(wc -l "${out}" | awk '{print $1}')
#cat_file "${out}"
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines, got ${cnt}" && exit 1
grep -v '.go$' "${out}" &>/dev/null || (echo "bad content" && exit 1)

# ===================================================
title ">>> test find hidden files <<<"
"${bin}" find -c "${catalog}" '**/.*' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../" -mindepth 1 -name '.*' | wc -l | awk '{print $1}')
#find "${cur}/../" -name '.*'
cnt=$(wc -l "${out}" | awk '{print $1}')
#cat_file "${out}"
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines, got ${cnt}" && exit 1
grep -v '.go$' "${out}" &>/dev/null || (echo "bad content" && exit 1)

# ===================================================
title ">>> test find .git files <<<"
"${bin}" find -c "${catalog}" --format=filename '**/.git/**' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find "${cur}/../.git/" -exec basename {} \; | wc -l | awk '{print $1}')
cnt=$(wc -l "${out}" | awk '{print $1}')
#cat_file "${out}"
#diff <(cat "${out}" | sort) <(find "${cur}/../.git/" -exec basename {} \; | sort)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines, got ${cnt}" && exit 1
grep -v '.go$' "${out}" &>/dev/null || (echo "bad content" && exit 1)

# index
catalog="${tmpd}/catalog2"
"${bin}" index -a -C -c "${catalog}" "${cur}/../internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ===================================================
title ">>> test find tree <<<"
"${bin}" find -c "${catalog}" 'tree' | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
expected=$(find internal -iname '*tree*' | wc -l | awk '{print $1}')
cnt=$(wc -l "${out}" | awk '{print $1}')
#cat_file "${out}"
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines, got ${cnt}" && exit 1

echo "test $(basename "${0}") OK!"
exit 0
