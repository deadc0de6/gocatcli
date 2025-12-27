#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test for nested directories with same names
# for https://github.com/deadc0de6/gocatcli/issues/3
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
tmpx=$(mktemp -d --suffix='-dotdrop-tests' || mktemp -d)
clear_on_exit "${tmpd}"
clear_on_exit "${tmpx}"

catalog="${tmpd}/catalog"
out="${tmpd}/output.txt"

# create hierarchy
top="${tmpx}/top"
sup="${top}/artist"
sub="${sup}/artist"
subf="${sub}/sub-file"
supf="${sup}/sup-file"
topf="${top}/top-file"
mkdir -p "${sub}"
echo "sub-file" > "${subf}"
echo "sup-file" > "${supf}"
echo "top-file" > "${topf}"

# index
echo ">>> index dir <<<"
"${bin}" index -a -C --debug -c "${catalog}" --ignore="\.git" "${top}" top
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# ls
echo ">>> ls catalog <<<"
"${bin}" -c "${catalog}" ls -a -r | sed -e 's/\x1b\[[0-9;]*m//g' | sed 's/[[:space:]]*$//' > "${out}"

cat > "${tmpx}/expected" << _EOF
storage top
artist
  artist
    sub-file
  sup-file
top-file
_EOF

if hash delta >/dev/null 2>&1; then
  delta "${tmpx}/expected" "${out}"
else
  diff -w --suppress-common-lines "${tmpx}/expected" "${out}"
fi

echo "test $(basename "${0}") OK!"
exit 0
