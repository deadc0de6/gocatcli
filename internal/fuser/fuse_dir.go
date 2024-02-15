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
	"gocatcli/internal/utils"
	iofs "io/fs"
	"os"
	"syscall"
	"time"

	"github.com/anacrolix/fuse"
	"github.com/anacrolix/fuse/fs"
)

// FuseDir a fuse directory
type FuseDir struct {
	theTree *tree.Tree
	current node.Node
	fs      *FS
}

func getNodeType(theNode node.Node) fuse.DirentType {
	switch theNode.GetType() {
	case node.FileTypeArchive:
		return fuse.DT_Dir
	case node.FileTypeArchived:
		return fuse.DT_File
	case node.FileTypeDir:
		return fuse.DT_Dir
	case node.FileTypeFile:
		return fuse.DT_File
	case node.FileTypeStorage:
		return fuse.DT_Dir
	}
	return fuse.DT_Unknown
}

func nodeToDirent(theNode node.Node) fuse.Dirent {
	dirent := fuse.Dirent{
		Type: getNodeType(theNode),
		Name: theNode.GetName(),
	}
	return dirent
}

func nodeToFuse(theNode node.Node, theTree *tree.Tree, filesys *FS) fs.Node {
	fuseType := getNodeType(theNode)
	if fuseType == fuse.DT_Dir {
		sub := &FuseDir{
			theTree: theTree,
			current: theNode,
			fs:      filesys,
		}
		return sub
	}

	sub := &FuseFile{
		theTree: theTree,
		current: theNode,
		fs:      filesys,
	}
	return sub
}

// Attr directory attributes
func (h *FuseDir) Attr(_ context.Context, a *fuse.Attr) error {
	if h == nil {
		return syscall.ENOENT
	}

	if h.fs.debugMode {
		line := fmt.Sprintf("dir attr of: %v", h.current)
		log.ToFile(logPath, line)
	}

	a.Size = 512

	if h == nil || h.current == nil {
		a.Inode = 1
		a.Mtime = time.Now()
		a.Ctime = time.Now()
		a.Mode = os.ModeDir | 0755
	} else {
		a.Inode = utils.HashString64(h.current.GetPath())
		a.Mtime = time.Unix(h.current.GetMAccess(), 0)
		a.Ctime = time.Unix(h.current.GetMAccess(), 0)
		mode := iofs.FileMode(utils.ModeStrToInt(h.current.GetMode()))
		log.ToFile(logPath, fmt.Sprintf("mode %s -> %v", h.current.GetMode(), mode))
		a.Mode = os.ModeDir | mode
	}

	return nil
}

// Lookup looks up directory
func (h *FuseDir) Lookup(_ context.Context, name string) (fs.Node, error) {
	if h.fs.debugMode {
		line := fmt.Sprintf("dir lookup of %v for name \"%s\"", h.current, name)
		log.ToFile(logPath, line)
	}

	if h.current == nil {
		// root
		storage := h.theTree.GetStorageByName(name)
		if storage == nil {
			return nil, syscall.ENOENT
		}
		sub := &FuseDir{
			theTree: h.theTree,
			current: storage,
			fs:      h.fs,
		}
		return sub, nil
	}

	// children
	entries := h.current.GetDirectChildren()
	if entries == nil {
		return nil, syscall.ENOENT
	}
	for childName, child := range entries {
		if childName == name {
			return nodeToFuse(child, h.theTree, h.fs), nil
		}
	}
	return nil, syscall.ENOENT
}

// ReadDirAll read dir content
func (h *FuseDir) ReadDirAll(_ context.Context) ([]fuse.Dirent, error) {
	var tops []fuse.Dirent

	if h.fs.debugMode {
		line := fmt.Sprintf("dir readdirall for %v", h.current)
		log.ToFile(logPath, line)
	}

	if h.current == nil {
		// root - list storages
		storages := h.theTree.GetStorages()
		for _, storage := range storages {
			dirent := nodeToDirent(storage)
			tops = append(tops, dirent)
		}
	} else {
		// children
		entries := h.current.GetDirectChildren()
		for _, child := range entries {
			dirent := nodeToDirent(child)
			tops = append(tops, dirent)
		}
	}

	return tops, nil
}
