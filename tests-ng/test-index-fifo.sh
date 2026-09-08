#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2026, deadc0de6
#
# test indexing a directory containing a fifo (named pipe)
# must not hang
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

src="${tmpd}/src"
catalog="${tmpd}/catalog"

mkdir "${src}"
echo "hello" > "${src}/hello.txt"
mkfifo "${src}/testpipe"
ln -s hello.txt "${src}/symlink"

# index (guarded by timeout so a regression fails instead of hanging)
title ">>> test index with fifo <<<"
if ! timeout 30 "${bin}" index -c "${catalog}" -C -a "${src}" src; then
  echo "index failed or hung on special file"
  exit 1
fi
[ ! -e "${catalog}" ] && echo "catalog not created" && exit 1

# the fifo must be indexed (metadata only, not opened)
grep -q "testpipe" "${catalog}" || { echo "fifo not found in catalog" && exit 1; }
grep -q 'prw-' "${catalog}" || { echo "fifo mode not found in catalog" && exit 1; }

# the regular file must keep its checksum
grep -q "hello.txt" "${catalog}" || { echo "regular file not found in catalog" && exit 1; }
grep -q '"md5": "[a-f0-9]\{32\}"' "${catalog}" || { echo "checksum not computed for regular file" && exit 1; }

echo "test $(basename "${0}") OK!"
exit 0
