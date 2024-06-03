/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package walker

import (
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
	"gocatcli/internal/utils"
	"gocatcli/internal/walker/archives"
	"io/fs"
	"path/filepath"
	"regexp"
)

// Walker a walker
type Walker struct {
	tree         *tree.Tree
	withChecksum bool
	withArchive  bool
	ignores      []*regexp.Regexp
	noMime       bool
}

// walk walks a dir - returns nb children and error if any
func (w *Walker) walk(storageID int, walkPath string, storagePath string, parent node.Node) (int64, error) {
	var cnt int64
	//log.Debugf("walking %s with parent %s", walkPath, parent.GetName())
	children := parent.GetDirectChildren()
	err := filepath.WalkDir(walkPath, func(pathUnderRoot string, dentry fs.DirEntry, err error) error {
		if err != nil {
			log.Error(err)
			// skipping
			return nil
		}

		if pathUnderRoot == walkPath {
			//log.Debugf("skipping \"%s\"", pathUnderRoot)
			// skipping
			return nil
		}
		if dentry == nil {
			log.Errorf("cannot index %s", walkPath)
			// skipping
			return nil
		}

		log.Debugf("indexing %s", pathUnderRoot)

		info, err := dentry.Info()
		if err != nil {
			log.Errorf("cannot index %s: %v", walkPath, err)
		}

		if w.mustIgnore(pathUnderRoot) {
			// skipping
			if info.IsDir() {
				log.Infof("ignoring directory \"%s\"...", pathUnderRoot)
				return filepath.SkipDir
			}
			log.Infof("ignoring \"%s\"", pathUnderRoot)
			return nil
		}

		// create or update child
		child, ok := children[info.Name()]
		if !ok {
			// create
			log.Debugf("node \"%s\" created", info.Name())
			fpath, err := filepath.Rel(storagePath, pathUnderRoot)
			if err != nil {
				return err
			}
			child = node.NewFileNode(storageID, fpath, info)
			parent.AddChild(child)
			//child = parent.GetDirectChildren()[info.Name()] // get the correct pointer
		} else {
			// update
			log.Debugf("updating node \"%s\"", info.Name())
			child.Update(info)
		}
		cnt++

		if child != nil {
			//log.Debugf("walker found path:\"%s\" (parent:\"%s\")", pathUnderRoot, parent.GetName())

			// handle directory
			if node.IsDir(child) {
				subcnt, err := w.walk(storageID, pathUnderRoot, storagePath, child)
				if err != nil {
					log.Error(err)
				}
				cnt += subcnt
				return filepath.SkipDir
			}

			// handle mime type
			if !w.noMime {
				child.Mime = getMime(pathUnderRoot)
			}

			// handle checksums
			if w.withChecksum {
				chk, err := utils.ChecksumFileContent(pathUnderRoot)
				if err != nil {
					log.Error(err)
				} else {
					log.Debugf("checksumming %s", pathUnderRoot)
					child.Checksum = chk
				}
			}

			// handle archives
			if w.withArchive && archives.IsArchive(pathUnderRoot) {
				log.Debugf("%s is archive", pathUnderRoot)
				processArchive(pathUnderRoot, storageID, storagePath, child)
			}
		}
		return nil
	})

	return cnt, err
}

func processArchive(path string, storageID int, storagePath string, child *node.FileNode) {
	//defer func() {
	//	r := recover()
	//	if r != nil {
	//		log.Errorf("archive indexing failed for %s", path)
	//	}
	//}()
	archived, _ := archives.GetFiles(path)
	for _, arc := range archived {
		fpath, err := filepath.Rel(storagePath, path)
		if err != nil {
			log.Errorf("archive indexing failed for %s: %v", path, err)
			return
		}
		sub := node.NewArchivedFileNode(storageID, fpath, arc.FileInfo, arc.Path)
		child.AddChild(sub)
	}
	child.Type = node.FileTypeArchive
}

// Walk walks the filesystem hierarchy
func (w *Walker) Walk(storageID int, walkPath string, storage *node.StorageNode) (int64, error) {
	cnt, err := w.walk(storageID, walkPath, walkPath, storage)

	type parentChild struct {
		parent node.Node
		child  node.Node
	}

	// clean unseen nodes
	var toRemove []*parentChild
	callback := func(n node.Node, _ int, parent node.Node) bool {
		if !n.Seen() {
			toRemove = append(toRemove, &parentChild{
				parent: parent,
				child:  n,
			})
			return false
		}
		return true
	}
	w.tree.ProcessChildren(storage, true, callback, -1)

	// we do in two steps since ProcessChildren loops the children
	// and RemoveChild alters the children slice
	for _, torm := range toRemove {
		torm.parent.RemoveChild(torm.child)
	}

	// re-traverse tree to set total size of directory
	log.Debugf("calculating total sizes...")
	storage.RecursiveFillSize()

	return cnt, err
}

func (w *Walker) mustIgnore(path string) bool {
	for _, patt := range w.ignores {
		matched := patt.MatchString(path)
		if matched {
			return matched
		}
	}
	return false
}

// NewWalker creates a new walker on path
func NewWalker(t *tree.Tree, withChecksum bool, withArchive bool, ignores []*regexp.Regexp, noMime bool) *Walker {
	w := Walker{
		tree:         t,
		withChecksum: withChecksum,
		withArchive:  withArchive,
		ignores:      ignores,
		noMime:       noMime,
	}
	return &w
}
