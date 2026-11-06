/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package fuser

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/deadc0de6/gocatcli/internal/log"
	"github.com/deadc0de6/gocatcli/internal/tree"

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
	defer func() {
		// recovering from a panic in c.Close: the anacrolix/fuse
		// library leaves backendState nil on linux which makes
		// Close panic - the connection is already torn down when
		// reaching this point anyway
		defer func() {
			_ = recover()
		}()
		err := c.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	// handle SIGINT/SIGTERM/SIGHUP: unmount so fs.Serve returns
	// and the deferred close is executed (leaves no stale mount)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sig)
	go func() {
		<-sig
		log.Info("unmounting, please wait...")
		err := fuse.Unmount(mountpoint)
		if err != nil {
			log.Errorf("could not unmount %s: %s", mountpoint, err)
			// unmount failed (mount busy?): force exit on next signal
			log.Info("press ctrl-c again to force exit")
			<-sig
			os.Exit(1)
		}
	}()

	myFS := &FS{
		mountPoint: mountpoint,
		theTree:    theTree,
		debugMode:  debug,
	}
	err = fs.Serve(c, myFS)
	return err
}
