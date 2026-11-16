#!/usr/bin/env bash
# author: Hadi Susanto Lesmana (https://github.com/hadi-susanto)
# Copyright (c) 2026, Hadi Susanto Lesmana
#
# test ls on paths containing glob metacharacters (brackets)
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
tmpx=$(mktemp -d --suffix='-gocatcli-tests' || mktemp -d)
clear_on_exit "${tmpd}"
clear_on_exit "${tmpx}"

catalog="${tmpd}/catalog"
out="${tmpd}/output.txt"

# create hierarchy
top="${tmpx}/top"
bd="${top}/Show (2020) [BD Rip]"
grp="${top}/[Group] Title"
mkdir -p "${bd}" "${grp}"
echo "ep01" > "${bd}/ep01.mkv"
echo "ep02" > "${grp}/ep02.mkv"
echo "top-file" > "${top}/top-file"

# index
# ===================================================
title ">>> index dir <<<"
"${bin}" index -a -C --debug -c "${catalog}" "${top}" top
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls dir with brackets at the end
# ===================================================
title ">>> ls dir with brackets at the end <<<"
"${bin}" -c "${catalog}" ls -a "top/Show (2020) [BD Rip]" | sed -e 's/\x1b\[[0-9;]*m//g' | sed 's/[[:space:]]*$//' > "${out}"
cat_file "${out}"

cat > "${tmpx}/expected" << _EOF
ep01.mkv
_EOF

diff -w "${tmpx}/expected" "${out}"

# ls dir with brackets at the start
# ===================================================
title ">>> ls dir with brackets at the start <<<"
"${bin}" -c "${catalog}" ls -a "top/[Group] Title" | sed -e 's/\x1b\[[0-9;]*m//g' | sed 's/[[:space:]]*$//' > "${out}"
cat_file "${out}"

cat > "${tmpx}/expected" << _EOF
ep02.mkv
_EOF

diff -w "${tmpx}/expected" "${out}"

echo "test $(basename "${0}") OK!"
exit 0
