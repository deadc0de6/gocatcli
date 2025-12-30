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

tmpd=$(mktemp -d --suffix='-gocatcli-tests' || mktemp -d)
clear_on_exit "${tmpd}"

catalog="${tmpd}/catalog"
out="${tmpd}/output.txt"

# index
# ===================================================
title ">>> test index <<<"
"${bin}" index -a -C -c "${catalog}" --debug --ignore='**/.git/**' --ignore='**/.git' "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
# ===================================================
title ">>> test index ls <<<"
"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#cat_file "${out}"
expected=$(find "${cur}/../" -mindepth 1 -name '.git' -prune -o -print | wc -l | awk '{print $1}')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

catalog="${tmpd}/catalog2"

# index
# ===================================================
title ">>> test index <<<"
"${bin}" index -a -C -c "${catalog}" "${cur}/../internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
# ===================================================
title ">>> test index ls (2) <<<"
"${bin}" --debug -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
expected=$(find "${cur}/../internal" -mindepth 1 -print | wc -l | awk '{print $1}')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

catalog="${tmpd}/catalog3"

# ===================================================
title ">>> test index with more ignore <<<"
"${bin}" index -a -C -c "${catalog}" --debug --ignore='**/.git/**' --ignore='**/.git' --ignore='**/*.go' --ignore='**/*.md' "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
expected=$(find "${cur}/../" -mindepth 1 -name '.git' -prune -o -name '*.md' -prune -o -name '*.go' -prune -o -print | wc -l | awk '{print $1}')
cat_file "${out}"
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

catalog="${tmpd}/catalog4"

# ===================================================
title ">>> test index with ignore hidden <<<"
"${bin}" index -a -C -c "${catalog}" --debug --ignore='**/.*' "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
#cat_file "${out}"
expected=$(find "${cur}/../" -mindepth 1 -name '.*' -prune -o -print | wc -l | awk '{print $1}')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

catalog="${tmpd}/catalog5"

# ===================================================
title ">>> test index 33 <<<"
"${bin}" index -a -C -c "${catalog}" --debug --ignore='**/.git*/**' --ignore='**/.git' "${cur}/../" gocatcli
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
# ===================================================
title ">>> test index ls <<<"
"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
# shellcheck disable=SC2126
cat_file "${out}"
expected=$(find "${cur}/../" -mindepth 1 -name '.git*' -prune -o -print | wc -l | awk '{print $1}')
cnt=$(tail -n +2 "${out}" | sed '/^$/d' | wc -l)
[ "${cnt}" != "${expected}" ] && echo "expecting ${expected} lines got ${cnt}" && exit 1

echo "test $(basename "${0}") OK!"
exit 0
