/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deadc0de6/gocatcli/internal/log"
	"github.com/deadc0de6/gocatcli/internal/node"

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

	createOptStart       string
	createOptWithArchive bool
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&createOptStart, "path", "p", "", "catalog start path")
	createCmd.PersistentFlags().BoolVar(&createOptWithArchive, "archive", false, "create archived files/dirs too")
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
		for _, top := range rootTree.GetStorages() {
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
				if len(n.GetDirectChildren()) > 0 {
					log.Debugf("mkdir for archive %s", p)
					err := os.MkdirAll(p, os.ModePerm)
					if err != nil {
						log.Error(err)
					}
				} else {
					log.Debugf("touch for archive %s", p)
					fd, err := os.Create(p)
					if err != nil {
						log.Error(err)
					}
					err = fd.Close()
					if err != nil {
						log.Error(err)
					}
				}
			case node.FileTypeArchived:
				if createOptWithArchive {
					sub := filepath.Dir(p)
					log.Debugf("mkdir for archived %s", sub)
					err := os.MkdirAll(sub, os.ModePerm)
					if err != nil {
						log.Error(err)
					}
					// handle files inside archive
					isArchivedDir := false
					if n.GetMode() == "" {
						// in doubt or when imported from catcli
						// check the children
						if len(n.GetDirectChildren()) > 0 {
							isArchivedDir = true
						}
					} else if node.IsModeDir(n) {
						isArchivedDir = true
					}

					if isArchivedDir {
						log.Debugf("mkdir for archived %s", p)
						err := os.MkdirAll(p, os.ModePerm)
						if err != nil {
							log.Error(err)
						}
					} else {
						log.Debugf("touch for archived %s", p)
						fd, err := os.Create(p)
						if err != nil {
							log.Error(err)
						}
						err = fd.Close()
						if err != nil {
							log.Error(err)
						}
					}
				}
			case node.FileTypeDir:
				log.Debugf("mkdir for dir %s", p)
				err := os.MkdirAll(p, os.ModePerm)
				if err != nil {
					log.Error(err)
				}
			case node.FileTypeFile:
				log.Debugf("touch for file %s", p)
				fd, err := os.Create(p)
				if err != nil {
					log.Error(err)
				}
				err = fd.Close()
				if err != nil {
					log.Error(err)
				}
			}
			return true
		}

		rootTree.ProcessChildren(n, true, callback, -1)
	}
	return nil
}
