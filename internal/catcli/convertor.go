/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package catcli

import (
	"encoding/json"
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
	"os"
	"time"
)

// readCatalog read catalog json file
func readCatalog(path string) (*Top, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := fd.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	var top Top
	err = json.NewDecoder(fd).Decode(&top)
	if err != nil {
		return nil, err
	}

	return &top, nil
}

func convertStorage(old *Storage) *node.StorageNode {
	meta := old.Attr
	newStorage := &node.StorageNode{
		ID:        node.DeriveStorageID(old.Name),
		Name:      old.Name,
		Size:      uint64(old.Size),
		Free:      uint64(old.Free),
		Total:     uint64(old.Total),
		IndexedAt: int64(old.TimeStamp),
		Type:      node.FileTypeStorage,
		Meta:      meta,
	}
	return newStorage
}

func convertFile(old *Node, storageID int) *node.FileNode {
	newFile := &node.FileNode{
		Name:      old.Name,
		RelPath:   old.RelPath,
		Checksum:  old.MD5,
		Type:      node.FileTypeFile,
		Size:      uint64(old.Size),
		Maccess:   int64(old.MAccess),
		StorageID: storageID,
	}
	return newFile
}

func convertDir(old *Node, storageID int) *node.FileNode {
	newDir := convertFile(old, storageID)
	newDir.Type = node.FileTypeDir
	return newDir
}

func convertArchive(old *Node, storageID int) *node.FileNode {
	newArchive := convertFile(old, storageID)
	newArchive.Type = node.FileTypeArchive
	return newArchive
}

func convertArchived(old *Node, storageID int) *node.FileNode {
	newArchived := convertFile(old, storageID)
	newArchived.Type = node.FileTypeArchived
	return newArchived
}

func debugNode(old *Node) {
	log.Debugf("name:%s archive:%s type:%s relpath:%s", old.Name, old.Archive, old.Type, old.RelPath)
}

func handleNode(storageID int, parent *node.FileNode, old *Node) (uint64, error) {
	var cnt uint64
	debugNode(old)
	switch old.Type {
	case nodeTypeFile:
		if len(old.Children) < 1 {
			// this is a file
			newFile := convertFile(old, storageID)
			parent.Children = append(parent.Children, newFile)
			cnt++
		} else {
			// this is an archive
			newArchive := convertArchive(old, storageID)
			parent.Children = append(parent.Children, newArchive)
			// recursively handle archive children
			for _, child := range old.Children {
				subcnt, err := handleNode(storageID, newArchive, child)
				if err != nil {
					return 0, err
				}
				cnt += subcnt
			}
		}

	case nodeTypeDir:
		// this is a directory
		newDir := convertDir(old, storageID)
		parent.Children = append(parent.Children, newDir)
		// recursively handle directory children
		for _, child := range old.Children {
			subcnt, err := handleNode(storageID, newDir, child)
			if err != nil {
				return 0, err
			}
			cnt += subcnt
		}
	case nodeTypeArchive:
		// this is an archived file (parent is the archive)
		newArchived := convertArchived(old, storageID)
		parent.Children = append(parent.Children, newArchived)
		// recursively handle archived children
		for _, child := range old.Children {
			subcnt, err := handleNode(storageID, newArchived, child)
			if err != nil {
				return 0, err
			}
			cnt += subcnt
		}
	default:
		return 0, fmt.Errorf("bad node type: %s", old.Type)
	}
	return cnt, nil
}

func handleStorage(_ *tree.Tree, old *Storage) (*node.StorageNode, error) {
	log.Debugf("storage: %#v", old)
	// create the storage in the newStorage tree
	newStorage := convertStorage(old)

	// handle storage children
	var total uint64
	for _, child := range old.Children {
		fake := &node.FileNode{}
		cnt, err := handleNode(newStorage.ID, fake, child)
		if err != nil {
			return nil, err
		}
		newStorage.Children = append(newStorage.Children, fake.Children...)
		total += cnt
	}

	newStorage.TotalFiles = total

	// re-traverse tree to set total size
	// and children count
	newStorage.RecursiveFillSize()

	return newStorage, nil
}

// Convert convert catalog to tree format
func Convert(version string, path string) (*tree.Tree, error) {
	top, err := readCatalog(path)
	if err != nil {
		return nil, err
	}

	log.Debugf("top name: %s", top.Name)
	log.Debugf("top type: %s", top.Type)
	log.Debugf("top nb children: %d", len(top.Children))

	t, err := tree.NewTree(version)
	if err != nil {
		return nil, err
	}
	t.Note = "converted from catcli"

	var topMeta *Meta
	for _, sub := range top.Children {
		// try to parse child
		var meta Meta
		err := json.Unmarshal(sub, &meta)
		if err == nil {
			// this is a meta
			if topMeta != nil {
				return nil, fmt.Errorf("two top meta node found")
			}
			topMeta = &meta
			continue
		}

		// match to storage
		var storage Storage
		err = json.Unmarshal(sub, &storage)
		if err != nil {
			return nil, err
		}

		sto, err := handleStorage(t, &storage)
		if err != nil {
			return nil, err
		}
		t.Storages = append(t.Storages, sto)
	}

	t.Created = topMeta.Meta.Created
	t.Updated = time.Now().Unix()

	return t, nil
}
