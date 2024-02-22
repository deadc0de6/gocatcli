/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package node

import (
	"fmt"
	"gocatcli/internal/utils"
	"io/fs"
	"sort"
	"strings"
	"time"
)

// GetName returns this node name
func (n *FileNode) GetName() string {
	return n.Name
}

// GetMAccess returns node modification date
func (n *FileNode) GetMAccess() int64 {
	return n.Maccess
}

// GetMode returns node mode
func (n *FileNode) GetMode() string {
	return n.Mode
}

// GetDirectChildren returns this node children
func (n *FileNode) GetDirectChildren() map[string]*FileNode {
	children := make(map[string]*FileNode, len(n.Children))
	for _, child := range n.Children {
		children[child.GetName()] = child
	}
	return children
}

// GetSortedDirectChildren returns children sorted by names
func (n *FileNode) GetSortedDirectChildren() []*FileNode {
	sort.Slice(n.Children, func(i, j int) bool {
		left := n.Children[i]
		right := n.Children[j]
		if sortByType {
			// first sort by type
			if IsDir(left) && !IsDir(right) {
				return false
			}
			if !IsDir(left) && IsDir(right) {
				return true
			}
		}

		// then by name
		leftName := strings.ToLower(left.GetName())
		rightName := strings.ToLower(right.GetName())
		return leftName < rightName
	})
	return n.Children
}

// GetPath returns the node relative path to its parent
func (n *FileNode) GetPath() string {
	return n.RelPath
}

// AddChild adds a child to this node
func (n *FileNode) AddChild(child *FileNode) {
	n.Children = append(n.Children, child)
}

// RemoveChild removes a child from this node
func (n *FileNode) RemoveChild(removeMe Node) {
	var newChildrenSlice []*FileNode
	for _, child := range n.Children {
		if child.GetName() != removeMe.GetName() {
			newChildrenSlice = append(newChildrenSlice, child)
		}
	}
	n.Children = newChildrenSlice
}

// GetType returns the node type
func (n *FileNode) GetType() FileType {
	return n.Type
}

// GetAttr returns the node attribute as string
func (n *FileNode) GetAttr(rawSize bool, long bool, extra bool) map[string]string {
	attrs := make(map[string]string)

	if !long {
		return attrs
	}

	// size
	size := fmt.Sprintf("%d", n.Size)
	if !rawSize {
		size = utils.SizeToHuman(n.Size)
	}
	attrs["size"] = size

	// maccess
	attrs["mode"] = n.Mode

	// type
	attrs["type"] = string(n.Type)

	// maccess
	tstr := utils.DateToString(n.Maccess)
	attrs["maccess"] = tstr

	if !extra {
		return attrs
	}

	// checksum
	if len(n.Checksum) > 0 {
		attrs["checksum"] = string(n.Checksum)
	}

	// index at
	indexed := utils.DateToString(n.IndexedAt)
	attrs["indexed"] = indexed

	// mime type
	mime := string(n.Mime)
	attrs["mime"] = mime

	// get extras
	if len(n.Extra) > 0 {
		extras := strings.Split(n.Extra, ",")
		for _, extra := range extras {
			fields := strings.Split(extra, ":")
			if len(fields) == 2 {
				attrs[fields[0]] = fields[1]
			}
		}
	}

	// nb children
	attrs["children"] = fmt.Sprint(len(n.Children))

	return attrs
}

// IsExec is file executable
func (n *FileNode) IsExec() bool {
	return strings.Count(n.Mode, "x") == 3
}

// GetSize returns this node size
func (n *FileNode) GetSize() uint64 {
	return n.Size
}

// SetSize sets the node size field
func (n *FileNode) SetSize(size uint64) {
	n.Size = size
}

// recursiveFillSize return size of subtree and cnt of file
func (n *FileNode) recursiveFillSize() (uint64, uint64) {
	if !MayHaveChildren(n) {
		return n.GetSize(), 1
	}

	// handle directory
	var size uint64
	var cnt uint64
	for _, child := range n.Children {
		childSize, childCnt := child.recursiveFillSize()
		size += childSize
		cnt += childCnt
	}

	n.Size = size
	return size, cnt
}

// Seen boolean to indicate if node was seen last update
func (n *FileNode) Seen() bool {
	return n.seen
}

// Update updates the node info
func (n *FileNode) Update(info fs.FileInfo) {
	n.Name = info.Name()
	n.Size = uint64(info.Size())
	n.Maccess = info.ModTime().Unix()
	n.IndexedAt = time.Now().Unix()
	n.Mode = info.Mode().String()
	n.seen = true
}

// NewArchivedFileNode creates a new archived file node
func NewArchivedFileNode(storageID int, path string, info fs.FileInfo, nameInsideArchive string) *FileNode {
	node := NewFileNode(storageID, path, info)
	if node == nil {
		return node
	}
	node.Type = FileTypeArchived
	node.Name = nameInsideArchive
	node.seen = true
	return node
}

// NewFileNode creates a new file node
func NewFileNode(storageID int, path string, info fs.FileInfo) *FileNode {
	if info == nil {
		return nil
	}
	typ := FileTypeFile
	if info.IsDir() {
		typ = FileTypeDir
	}
	node := FileNode{
		Type:      FileType(typ),
		StorageID: storageID,
	}
	node.Update(info)
	node.seen = true
	node.RelPath = path

	return &node
}
