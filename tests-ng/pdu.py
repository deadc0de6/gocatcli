#!/usr/bin/env python3
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#
# a naive python du alternative
# ./tests-ng/pdu.py . --ignore='*/.git*' -H
#

import os
import argparse
import fnmatch
from typing import List


NAME = 'pdu'
DEBUG = True
IGNORE = [
    '*/.DS_Store',
]


def debug(txt: str):
    """debug output"""
    if not DEBUG:
        return
    print(f'[DEBUG] {txt}')


def size_to_str(size: int, human: bool = False, digits=None) -> str:
    """size to string"""
    div = 1024.
    suf = ['B', 'K', 'M', 'G', 'T', 'P']
    if not human or size < div:
        # not human
        return f'{size}'
    sz = float(size)
    for i in suf:
        if sz < div:
            return f'{round(sz, digits)}{i}'
            # return f'{int(sz)}{i}'
        sz = sz / div
    sufix = suf[-1]
    return f'{round(size, digits)}{sufix}'
    # return f'{int(sz)}{sufix}'


def must_ignore(path: str, patterns: List[str]) -> bool:
    """must ignore path"""
    if not patterns:
        return False
    lst = [fnmatch.fnmatch(path, patt) for patt in patterns]
    debug(f'{path} -> {lst}')
    return any(lst)


def main(path: str, human: bool,
         ign: List[str] = []):
    """entry point"""
    if not os.path.exists(path):
        print(f'[ERROR] {path} does not exist')
        return False
    if not ign:
        ign = []
    ign.extend(IGNORE)
    total = 0
    for root, _, files in os.walk(path, topdown=True):
        dirsz = 0
        if must_ignore(root, ign):
            debug(f'ignore root {root}')
            continue
        files.sort(key=lambda f: f[0] == '.')
        for file in files:
            fpath = os.path.join(root, file)
            if must_ignore(fpath, ign):
                debug(f'ignore sub {fpath}')
                continue
            # size = os.path.getsize(fpath)
            stat_info = os.stat(fpath)
            size = stat_info.st_size
            dirsz += size
            total += size
        if root != path and dirsz > 0:
            print_line(dirsz, get_norm_path(root, path), human, digits=1)
    print_line(total, get_norm_path(path, path), human)
    return True

def get_norm_path(path: str, base: str) -> str:
    npath = os.path.normpath(path)
    bnpath = os.path.normpath(base)
    bbasename = os.path.basename(bnpath)
    rel = os.path.relpath(npath, bnpath)
    if rel == '.':
        return bbasename
    return os.path.join(bbasename, rel)

def print_line(size, path, human, digits=None):
    hsize = size_to_str(size, human=human, digits=digits)
    if not human:
        print(f'{hsize:<10}{path}')
        return
    print(f'{hsize:<10}{path}')

if __name__ == "__main__":
    parser = argparse.ArgumentParser(prog=NAME,
                                     description='python du')
    parser.add_argument('path')
    parser.add_argument('-H', '--human',
                        action='store_true')
    parser.add_argument('-i', '--ignore',
                        nargs='+')
    parser.add_argument('-d', '--debug',
                        action='store_true')
    args = parser.parse_args()
    DEBUG = args.debug
    main(args.path, args.human, ign=args.ignore)
