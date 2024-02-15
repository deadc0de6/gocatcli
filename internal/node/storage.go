/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package node

import (
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/utils"
	"sort"
	"strings"
	"time"
)

// GetName returns this node name
func (n *StorageNode) GetName() string {
	return n.Name
}

// GetMAccess returns node modification date
func (n *StorageNode) GetMAccess() int64 {
	return n.IndexedAt
}

// GetMode returns the storage mode
func (n *StorageNode) GetMode() string {
	return "drwxr-xr-x" // 0755
}

// SetMeta returns this node name
func (n *StorageNode) SetMeta(meta string) {
	n.Meta = meta
}

// IsDir returns true if can enter
func (n *StorageNode) IsDir() bool {
	return true
}

// GetDirectChildren returns this node children
func (n *StorageNode) GetDirectChildren() map[string]*FileNode {
	if n == nil || n.Children == nil {
		return nil
	}
	children := make(map[string]*FileNode, len(n.Children))
	for _, child := range n.Children {
		children[child.GetName()] = child
	}
	return children
}

// GetSortedDirectChildren returns children sorted by names
func (n *StorageNode) GetSortedDirectChildren() []*FileNode {
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

// GetPath returns the node relative path
func (n *StorageNode) GetPath() string {
	return n.Path
}

// AddChild adds a new child to this node
func (n *StorageNode) AddChild(child *FileNode) {
	n.Children = append(n.Children, child)
}

// Seen boolean to indicate if node was seen last update
func (n *StorageNode) Seen() bool {
	return true
}

// RecursiveFillSize recusively fills each node total size
func (n *StorageNode) RecursiveFillSize() {
	var size uint64
	var cnt uint64

	for _, child := range n.Children {
		subsize, subcnt := child.recursiveFillSize()
		size += subsize
		cnt += subcnt
	}

	log.Debugf("setting storage total size to %dB and total files to %d", size, cnt)
	n.SetSize(size)
	n.TotalFiles = cnt
}

// RemoveChild removes a child from this node
func (n *StorageNode) RemoveChild(removeMe Node) {
	var newChildrenSlice []*FileNode
	for _, child := range n.Children {
		if child.GetName() != removeMe.GetName() {
			newChildrenSlice = append(newChildrenSlice, child)
		}
	}
	n.Children = newChildrenSlice
}

// GetType returns the node type
func (n *StorageNode) GetType() FileType {
	return n.Type
}

func sizeToString(sz uint64, rawSize bool) string {
	if rawSize {
		return fmt.Sprintf("%d", sz)
	}
	return utils.SizeToHuman(sz)
}

// GetAttr returns the node attribute as string
func (n *StorageNode) GetAttr(rawSize bool, long bool) map[string]string {
	attrs := make(map[string]string)
	attrs["nbfiles"] = fmt.Sprintf("%d", n.TotalFiles)
	attrs["size"] = sizeToString(n.Size, rawSize)
	total := sizeToString(n.Total, rawSize)
	attrs["fs_size"] = total
	freePercent := "??"
	if n.Total != 0 {
		freePercent = fmt.Sprintf("%d%%", n.Free*100/n.Total)
	}
	attrs["fs_free"] = freePercent
	used := sizeToString(n.Total-n.Free, rawSize)
	attrs["fs_du"] = fmt.Sprintf("%s/%s", used, total)
	attrs["indexed"] = utils.DateToString(n.IndexedAt)

	if !long {
		return attrs
	}

	attrs["meta"] = n.Meta
	tags := n.Tags
	sort.Strings(tags)
	attrs["tags"] = strings.Join(tags, ",")

	return attrs
}

// Tag adds a tag to a storage
func (n *StorageNode) Tag(tag string) {
	for _, t := range n.Tags {
		if t == tag {
			return
		}
	}
	n.Tags = utils.UniqStrings(n.Tags, []string{tag})
}

// Untag removes a tag from storage
func (n *StorageNode) Untag(tag string) {
	var slice []string
	for _, t := range n.Tags {
		if t != tag {
			slice = append(slice, t)
		}
	}
	n.Tags = utils.UniqStrings(slice, []string{})
}

// GetSize returns this node size
func (n *StorageNode) GetSize() uint64 {
	return n.Size
}

// SetSize sets the node size field
func (n *StorageNode) SetSize(size uint64) {
	n.Size = size
}

// UpdateStorage updates a storage fields
func (n *StorageNode) UpdateStorage(fsPath string, path string, meta string, tags []string) {
	free, total := utils.DiskUsage(fsPath)
	n.Tags = utils.UniqStrings(n.Tags, tags)
	n.Free = free
	n.Total = total
	n.Path = path
	n.Meta = meta
	n.IndexedAt = time.Now().Unix()
}

// DeriveStorageID derive id from storage name
func DeriveStorageID(name string) int {
	now := time.Now().Format("2006-01-02 15:04:05")
	return utils.HashString(name + now)
}

// NewStorageNode creates a new storage node
func NewStorageNode(name string, fsPath string, path string, meta string, tags []string) *StorageNode {

	storage := StorageNode{
		ID:   DeriveStorageID(name),
		Name: name,
		Type: FileTypeStorage,
		Meta: meta,
	}
	storage.UpdateStorage(fsPath, path, meta, tags)
	return &storage
}
