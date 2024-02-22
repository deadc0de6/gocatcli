/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package node

const (
	sortByType = false
)

// FileType node file type
type FileType string

// Node generic node interface
type Node interface {
	GetName() string
	GetDirectChildren() map[string]*FileNode
	GetSortedDirectChildren() []*FileNode
	GetPath() string
	GetType() FileType
	GetMAccess() int64
	GetMode() string
	GetAttr(bool, bool) map[string]string // rawsize, long
	GetSize() uint64
	SetSize(uint64)
	Seen() bool
	AddChild(*FileNode)
	RemoveChild(Node)
}

// MayHaveChildren returns true if the node may have children
func MayHaveChildren(n Node) bool {
	return n.GetType() == FileTypeDir || n.GetType() == FileTypeStorage
}

// IsDir returns true if node is a directory
func IsDir(n Node) bool {
	return n.GetType() == FileTypeDir
}

// IsStorage returns true if node is storage
func IsStorage(n Node) bool {
	return n.GetType() == FileTypeStorage
}
