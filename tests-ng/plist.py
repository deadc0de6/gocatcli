#!/usr/bin/env python3
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# a naive python find files alternative
# ./tests-ng/pdu.py . --ignore='*/.git*' -H
#


import os
import argparse
import fnmatch
from typing import List


NAME = 'plist'
DEBUG = True
IGNORE = [
    '*/.DS_Store',
]


def debug(txt: str):
    """debug output"""
    if not DEBUG:
        return
    print(f'[DEBUG] {txt}')


def must_ignore(path: str, patterns: List[str]) -> bool:
    """must ignore path"""
    if not patterns:
        return False
    lst = [fnmatch.fnmatch(path, patt) for patt in patterns]
    return any(lst)


def main(path: str, ign: List[str] = []):
    """entry point"""
    if not os.path.exists(path):
        print(f'[ERROR] {path} does not exist')
        return False
    if not ign:
        ign = []
    ign.extend(IGNORE)
    cnt = 0
    for root, dirs, files in os.walk(path, topdown=True):
        if must_ignore(root, ign):
            debug(f'ignore root {root}')
            continue
        for file in files:
            fpath = os.path.join(root, file)
            if must_ignore(fpath, ign):
                debug(f'ignore sub {fpath}')
                continue
            debug(f'file: {file}')
            cnt += 1
        for d in dirs:
            fpath = os.path.join(root, d)
            if must_ignore(fpath, ign):
                debug(f'ignore sub {fpath}')
                continue
            debug(f'dir: {d}')
            cnt += 1
    print(f'{cnt}')
    return True


if __name__ == "__main__":
    parser = argparse.ArgumentParser(prog=NAME,
                                     description='python du')
    parser.add_argument('path')
    parser.add_argument('-i', '--ignore',
                        nargs='+')
    parser.add_argument('-d', '--debug',
                        action='store_true')
    args = parser.parse_args()
    DEBUG = args.debug
    main(args.path, ign=args.ignore)
