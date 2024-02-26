/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package fuser

import (
	"gocatcli/internal/tree"

	"github.com/anacrolix/fuse"
	"github.com/anacrolix/fuse/fs"
)

var (
	logPath = "/tmp/gocatcli-fuser.log"
)

// FS fuse filesystem
type FS struct {
	mountPoint string
	theTree    *tree.Tree
	root       *FuseDir
	debugMode  bool
}

// Root returns the root dir handle
func (f *FS) Root() (fs.Node, error) {
	rd := &FuseDir{
		fs:      f,
		theTree: f.theTree,
	}
	f.root = rd
	return rd, nil
}

// Mount mount the tree
func Mount(theTree *tree.Tree, mountpoint string, debug bool) error {
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("gocatcli"),
		fuse.Subtype("gocatcli"),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	myFS := &FS{
		mountPoint: mountpoint,
		theTree:    theTree,
		debugMode:  debug,
	}
	err = fs.Serve(c, myFS)
	return err
}
