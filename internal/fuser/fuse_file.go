/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package fuser

import (
	"context"
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
	"gocatcli/internal/utilities"
	"io/fs"
	"time"

	"github.com/anacrolix/fuse"
)

// FuseFile a file in fuse filesystem
type FuseFile struct {
	theTree *tree.Tree
	current node.Node
	fs      *FS
}

// Attr file attribute
func (h *FuseFile) Attr(_ context.Context, a *fuse.Attr) error {
	if h.fs.debugMode {
		line := fmt.Sprintf("%v file attr", h.current)
		log.ToFile(logPath, line)
	}

	a.Inode = utilities.HashString64(h.current.GetPath())
	a.Mode = 0755
	a.Size = h.current.GetSize()
	a.Atime = time.Unix(h.current.GetMAccess(), 0)
	a.Mtime = time.Unix(h.current.GetMAccess(), 0)
	a.Mode = fs.FileMode(utilities.ModeStrToInt(h.current.GetMode()))
	log.ToFile(logPath, fmt.Sprintf("mode %s -> %v", h.current.GetMode(), a.Mode))
	return nil
}
