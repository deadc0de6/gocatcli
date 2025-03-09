/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package tree

import (
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"gocatcli/internal/utils"
	"path/filepath"
	"strings"
	"time"
)

const (
	toolName = "gocatcli - https://github.com/deadc0de6/gocatcli"
)

// Tree the tree
type Tree struct {
	Storages []*node.StorageNode `json:"storages" toml:"storages"`
	Tool     string              `json:"tool" toml:"tool"`
	Version  string              `json:"version" toml:"version"`
	Created  int64               `json:"created" toml:"created"`
	Updated  int64               `json:"updated" toml:"updated"`
	Note     string              `json:"note" toml:"note"`
	//Nodes    map[string]*node.FileNode `json:"-" toml:"-"`
}

// ProcessCallback will be called with the current node, its depth and its parent
type ProcessCallback func(current node.Node, depth int, parent node.Node) bool

func matchPath(current node.Node, path string) bool {
	name := current.GetName()
	matched, err := filepath.Match(path, name)
	if !matched || err != nil {
		return false
	}
	return true
}

// recursively find nodes matching path entries
func (t *Tree) descendNodeWithPath(current node.Node, paths []string) []node.Node {
	if len(paths) < 1 {
		return nil
	}

	if !matchPath(current, paths[0]) {
		return nil
	}
	log.Debugf("descendNodeWithPath \"%s\" match pattern \"%s\"", current.GetName(), paths[0])

	if node.IsDir(current) && len(paths) > 1 {
		// descent
		var subs []node.Node
		for _, child := range current.GetDirectChildren() {
			matchedSubs := t.descendNodeWithPath(child, paths[1:])
			if matchedSubs != nil {
				subs = append(subs, matchedSubs...)
			}
		}
		return subs
	}

	return []node.Node{current}
}

// GetNodesFromPath returns all nodes which path match the path argument, nil otherwise
// if storage is defined, will only look into its nodes
// The path argument can use regexp/wildcards
func (t *Tree) GetNodesFromPath(path string) []node.Node {
	log.Debugf("GetNodesFromPath for path \"%s\"", path)
	if len(path) < 1 {
		// no path provided
		return nil
	}

	// split path for search
	paths := utils.SplitPath(path)
	if len(paths) < 1 {
		return nil
	}

	// find the storage nodes matching
	tops := []node.Node{}
	for _, top := range t.GetStorages() {
		if matchPath(top, paths[0]) {
			log.Debugf("selected top: %s", top.GetName())
			tops = append(tops, top)
		}
	}

	var found []node.Node
	if len(paths) > 1 {
		for _, top := range tops {
			for _, child := range top.GetDirectChildren() {
				sub := t.descendNodeWithPath(child, paths[1:])
				if sub != nil {
					found = append(found, sub...)
				}
			}
		}
	}

	if len(found) < 1 {
		return tops
	}
	return found
}

// GetStorages returns all storage for this tree
func (t *Tree) GetStorages() []*node.StorageNode {
	return t.Storages
}

// ProcessChildren process the tree depth-first
// when callback is defined, return is empty, when undefined, return is filled
// callback arguments: current node, depth, parent node
// maxDepth: max depth to process (-1 for infinite)
// callback returns true to continue down the tree or false to stop
func (t *Tree) ProcessChildren(start node.Node, hiddenToo bool, callback ProcessCallback, maxDepth int) []node.Node {
	log.Debugf("processing nodes children from \"%s\" (show hidden:%v)", start.GetName(), hiddenToo)
	return t.processChildren(start, callback, 0, maxDepth, hiddenToo)
}

// ProcessChildren process the tree depth-first
// when callback is defined, return is empty, when undefined, return is filled
// callback must return a boolean, true to continue processing the children
// callback arguments: current node, depth, parent node
// maxDepth: max depth to process (-1 for infinite)
func (t *Tree) processChildren(n node.Node, callback ProcessCallback, indent int, maxDepth int, hiddenToo bool) []node.Node {
	var retNodes []node.Node
	if n == nil {
		return nil
	}

	if maxDepth > -1 && indent > maxDepth {
		return nil
	}

	children := n.GetSortedDirectChildren()
	for _, child := range children {
		// filter hidden files
		if strings.HasPrefix(child.GetName(), ".") && !hiddenToo {
			continue
		}
		// callback
		if callback != nil {
			// callback
			if !callback(child, indent, n) {
				// callback asked us to stop
				log.Debugf("callback asked to stop at %s", child.GetName())
				continue
			}
		} else {
			retNodes = append(retNodes, child)
		}

		// recursive process child
		subs := t.processChildren(child, callback, indent+1, maxDepth, hiddenToo)
		retNodes = append(retNodes, subs...)
	}
	return retNodes
}

// GetStorageNode returns the storage for this node
func (t *Tree) GetStorageNode(n node.Node) *node.StorageNode {
	if n.GetType() == node.FileTypeStorage {
		sto := n.(*node.StorageNode)
		return sto
	}
	fnode := n.(*node.FileNode)
	id := fnode.StorageID
	return t.GetStorageByID(id)
}

// GetStorageByID returns the storage by id
func (t *Tree) GetStorageByID(id int) *node.StorageNode {
	for _, storage := range t.Storages {
		if storage.ID == id {
			return storage
		}
	}
	return nil
}

// GetStorageByName returns the storage by name
func (t *Tree) GetStorageByName(name string) *node.StorageNode {
	for _, storage := range t.Storages {
		if storage.GetName() == name {
			return storage
		}
	}
	return nil
}

// RemoveStorage removes a storage
func (t *Tree) RemoveStorage(name string) {
	var newStorage []*node.StorageNode
	for _, storage := range t.Storages {
		if storage.GetName() != name {
			newStorage = append(newStorage, storage)
		}
	}
	t.Storages = newStorage
}

// NewTree creates a new tree
func NewTree(version string) (*Tree, error) {
	tree := Tree{
		Version: version,
		Created: time.Now().Unix(),
		Tool:    toolName,
	}
	return &tree, nil
}
