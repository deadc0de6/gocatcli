#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# test create command
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
create_dir="${tmpd}/create"

# index
"${bin}" --debug index -a -C -c "${catalog}" "${cur}/../internal" internal
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

echo ">>> create <<<"
"${bin}" --debug create -c "${catalog}" "${create_dir}" | sed -e 's/\x1b\[[0-9;]*m//g' > "${out}"
tree "${create_dir}"
diff "${cur}/../internal/" "${create_dir}/internal/"

echo "test $(basename "${0}") OK!"
exit 0
