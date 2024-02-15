/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:    "create <local-path>",
		Short:  "Create filesystem hierarchy locally",
		PreRun: preRun(true),
		Args:   cobra.ExactArgs(1),
		RunE:   create,
	}

	createOptStart string
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&createOptStart, "path", "p", "", "catalog start path")
}

func create(_ *cobra.Command, args []string) error {
	localPath := args[0]

	// get the base paths for start
	var startNodes []node.Node
	if len(createOptStart) > 0 {
		startNodes = getStartPaths(createOptStart)
		if startNodes == nil {
			return fmt.Errorf("no such start path: \"%s\"", createOptStart)
		}
	} else {
		for _, top := range loadedTree.GetStorages() {
			startNodes = append(startNodes, top)
		}
	}

	// create base dir
	err := os.MkdirAll(localPath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, n := range startNodes {
		subPath := filepath.Join(localPath, n.GetName())
		log.Debugf("creating %s under %s", n.GetName(), subPath)

		// create base dir
		err := os.MkdirAll(subPath, os.ModePerm)
		if err != nil {
			return err
		}

		// travers the tree from a node and create hierarchy locally
		callback := func(n node.Node, _ int, _ node.Node) bool {
			p := filepath.Join(subPath, n.GetPath())
			switch n.GetType() {
			case node.FileTypeArchive:
				log.Debugf("mkdir %s", p)
				err := os.MkdirAll(p, os.ModePerm)
				if err != nil {
					log.Error(err)
				}
			case node.FileTypeArchived:
				log.Debugf("touch %s", p)
				fd, err := os.Create(p)
				if err != nil {
					log.Error(err)
				}
				fd.Close()
			case node.FileTypeDir:
				log.Debugf("mkdir %s", p)
				err := os.MkdirAll(p, os.ModePerm)
				if err != nil {
					log.Error(err)
				}
			case node.FileTypeFile:
				log.Debugf("touch %s", p)
				fd, err := os.Create(p)
				if err != nil {
					log.Error(err)
				}
				fd.Close()
			}
			return true
		}

		loadedTree.ProcessChildren(n, true, callback, -1)
	}
	return nil
}
