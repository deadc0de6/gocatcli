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


def debug(txt: str):
    """debug output"""
    if not DEBUG:
        return
    print(f'[DEBUG] {txt}')


def size_to_str(size: float, human: bool = False) -> str:
    """size to string"""
    div = 1024.
    suf = ['B', 'K', 'M', 'G', 'T', 'P']
    if not human or size < div:
        return f'{size}'
    for i in suf:
        if size < div:
            return f'{size:.1f}{i}'
        size = size / div
    sufix = suf[-1]
    return f'{size:.1f}{sufix}'


def must_ignore(path: str, patterns: List[str]) -> bool:
    """must ignore path"""
    if not patterns:
        return False
    lst = [fnmatch.fnmatch(path, patt) for patt in patterns]
    debug(f'{path} -> {lst}')
    return any(lst)


def main(path: str, human: bool,
         ign: List[str] = None):
    """entry point"""
    if not os.path.exists(path):
        print(f'[ERROR] {path} does not exist')
        return False
    total = 0
    for root, _, files in os.walk(path, topdown=True):
        dirsz = 0
        if must_ignore(root, ign):
            debug(f'ignore root {root}')
            continue
        for file in files:
            fpath = os.path.join(root, file)
            if must_ignore(fpath, ign):
                debug(f'ignore sub {fpath}')
                continue
            size = os.path.getsize(fpath)
            dirsz += size
            total += size
        if root != path:
            print(f'{size_to_str(dirsz, human=human)} {root}')
    print(f'{size_to_str(total, human=human)} {path}')
    return True


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
