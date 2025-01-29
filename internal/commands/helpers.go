/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"path/filepath"
	"strings"
)

func getStartPaths(path string) []node.Node {
	log.Debugf("getting start paths from \"%s\"", path)
	storages := rootTree.GetStorages()
	// do not mess with pattern
	if !strings.Contains(path, "*") && len(storages) == 1 {
		// complete if single storage
		name := storages[0].GetName()
		if !strings.HasPrefix(path, name) {
			path = filepath.Join(name, path)
		}
	}
	return rootTree.GetNodesFromPath(path)
}
