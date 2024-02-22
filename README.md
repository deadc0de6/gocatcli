# gocatcli

[![Tests Status](https://github.com/deadc0de6/gocatcli/workflows/tests/badge.svg)](https://github.com/deadc0de6/gocatcli/actions)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[![Donate](https://img.shields.io/badge/donate-KoFi-blue.svg)](https://ko-fi.com/deadc0de6)

*[gocatcli](https://github.com/deadc0de6/gocatcli) is a catalog tool for your offline data*

Did you ever wanted to find back that specific file that should be on one of your backup DVDs or one of your external hard drives?
You usually go through all of them hoping to find the right one on the first try?
`gocatcli` indexes external media in a catalog file and allows to quickly find
specific files or even navigate in the catalog as if it was a mounted drive.

Features:

* Index filesystem/directories in a catalog
* Multiple ways to explore the catalog (`ls`, `find`, `tree`, tui browser, `fzf`, etc)
* Use wildcards to search for files
* archives support (index their content as well)
* Save catalog to json for easy versioning with git
* Multiple outputs (`csv`, etc)
* Mount file using fuse
* Re-create locally the catalog hierarchy
* ... and much more

Quick start:
```bash
gocatcli index <some-path>
gocatli nav
gocatli tree
```

# Why gocatcli?

`gocatcli` gives the ability to navigate, explore and find your files that are stored on external media
(DVDs, hard drives, USB sticks, etc) when those are not connected.
`gocatcli` can just as easily index any arbitrary directories.

---

**Table of Contents**

* [Installation](#installation)
* [Usage](#usage)

  * [Index data](#index-data)
  * [Reindex and update](#reindex-and-update)
  * [Index archives](#index-archives-and-their-content)
  * [Navigate with ls](#navigate-with-ls)
  * [File browser](#file-browser)
  * [Tree view](#tree-view)
  * [Find files](#find-files)
  * [Find files with fzf](#find-files-with-fzf)
  * [Disk usage](#disk-usage)
  * [Create hierarchy locally](#create-hierarchy-locally)
  * [Mount the catalog filesystem](#mount-filesystem)
  * [Edit storage](#edit-storage)
  * [Output formats](#output-formats)
  * [Convert catcli catalog](#convert-catcli-catalog)

* [From catcli to gocatcli](#from-catcli-to-gocatcli)
* [Contribution](#contribution)
* [Thank you](#thank-you)

# Installation

Pick a binary in the [releases](https://github.com/deadc0de6/gocatcli/releases)

Or if you have go installed, you can installed it directly
```bash
## You need at least golang 1.22
$ go install -v github.com/deadc0de6/gocatcli/cmd/gocatcli@latest
$ gocatcli --help
```

Or if you want to compile it yourself
```bash
## You need at least golang 1.22
$ go mod tidy
$ make
$ ./bin/gocatcli --help
```

# Usage

The primary use of gocatcli is to index your data (external hardrives, etc) into a library
and then have the ability to browse their content (as well as search, navigate, etc) while these
are stored away.

```bash
$ gocatcli --help
$ gocatcli index --help
$ gocatcli ls --help
$ gocatcli tree --help
$ gocatcli nav --help
...
```

Wildcards are supported in the `<path>` arguments of all commands and provide a way
to explore the catalog using something like `'storage/directory*/sub-directory*'`,
Make sure to quote these on the command line to avoid your shell interpreting the
wildcards.

All command line switches can be provided using environment variables by
prefixing with `GOCATCLI_` and adding the switch name in capital and `-`
replaced with `-`. For example `--catalog` would be `GOCATCLI_CATALOG=mycatalog.catalog`.

You can generate shell completion scripts using `gocatcli completion`
```bash
## for zsh
## copy the file to a path within ${fpath}
$ gocatcli completion zsh
## for bash
## source the file
$ gocatcli completion bash
## for fish
## mkdir -p ~/.config/fish/completions
## cp gocatcli.fish ~/.config/fish/completions/gocatcli.fish
$ gocatcli completion fish
```

## Index data

```bash
$ gocatli index --help
```

Index any directories with
```bash
$ gocatcli index /some/directory
```

The below example ignores any file ending with `.go` or `.md` and anything in the `.git` directory:
```bash
$ gocatcli index ../gocatcli --ignore="*.go" --ignore="*.md" --ignore="*.git/*"
```

## Reindex and update

To re-index the content of an already indexed storage, simply re-run `index` on it
```bash
## initial indexing
$ ./bin/gocatcli index /tmp tmp-dir
## re-indexing...
$ ./bin/gocatcli index /tmp tmp-dir
A storage with the name "tmp-dir" already exists, update it? [y/N]: y
```

## Index archives and their content

`gocatcli` is able to index the content of archives.
Simply provide `-a --archive` to your `index` command.

Supported archive formats (from <https://github.com/mholt/archiver>):

* brotli (.br)
* bzip2 (.bz2)
* flate (.zip)
* gzip (.gz)
* lz4 (.lz4)
* snappy (.sz)
* xz (.xz)
* zlib (.zz)
* zstandard (.zst)
* .zip
* .tar (including any compressed variants like .tar.gz)
* .rar
* .7z

## Navigate with ls

```bash
$ gocatli ls --help
```

```bash
$ gocatcli ls
$ gocatcli ls storage-name/some/path
$ gocatcli ls 'storage-name/direc*/subdire*'
```

## File browser

A terminal file browser for your catalog
```bash
$ gocatli nav --help
```

## Tree view

```bash
$ gocatli tree --help
```

## Find files

With the `find` command you can easily find all files matching a specific
pattern. You can also limit your search to a specific path

```bash
$ gocatli find --help
## lists all files
$ gocatcli find
## find using pattern
$ gocatcli find pattern
## find using pattern and limit to a specific path
$ gocatcli find pattern -p some/path
$ gocatcli find pattern -p 'some/p*th'
```

## Find files with fzf

A terminal fzf file browser for your catalog
```bash
$ gocatli fzfind --help
```

## Disk usage

```bash
$ gocatli du --help
```

## Create hierarchy locally

```bash
$ gocatli create --help
```

## Mount filesystem

```bash
$ gocatli mount --help
```

## Edit storage

```bash
$ gocatli storage --help
```

List all storages with the `storage list` command.

Remove a storage and its children with the `storage rm` command.

Following commands allow to edit a storage and its fields:

* `meta`: storage description
* `tag`: add a tag to the storage
* `untag`: remove a tag from the storage

## Output formats

* `native`: ls-like output
* `csv`: csv
* `csv-with-header`: csv with header
* `tree`: tree
* `script`: generates a script to handle matches
* `debug`: debug output

## Convert catcli catalog

```bash
$ gocatli convert --help
```

# From catcli to gocatcli

<https://github.com/deadc0de6/catcli/>

```bash
$ gocatli convert --help
```

# Contribution

If you are having trouble installing or using `gocatcli`, open an issue.

If you want to contribute, feel free to do a PR.

The `tests.sh` script handles the linting and runs the tests.

# Thank you

If you like `gocatcli`, [buy me a coffee](https://ko-fi.com/deadc0de6).

# License

This project is licensed under the terms of the GPLv3 license.

