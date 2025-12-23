/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/deadc0de6/gocatcli/internal/log"
	"github.com/deadc0de6/gocatcli/internal/node"
	"github.com/deadc0de6/gocatcli/internal/tree"
	"github.com/deadc0de6/gocatcli/internal/utilities"
	"github.com/deadc0de6/gocatcli/internal/walker"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	indexCmd = &cobra.Command{
		Use:    "index <path> [<name>]",
		Short:  "Index a directory in the catalog",
		PreRun: preRun(false),
		Args:   cobra.RangeArgs(1, 2),
		RunE:   index,
	}

	indexOptTags     []string
	indexOptChecksum bool
	indexOptMeta     string
	indexOptArchive  bool
	indexOptIgnores  []string
	indexOptIndent   bool
	indexOptForce    bool
	indexOptNoMIME   bool
)

func init() {
	rootCmd.AddCommand(indexCmd)

	indexCmd.PersistentFlags().StringSliceVarP(&indexOptTags, "tag", "t", nil, "add a tag")
	indexCmd.PersistentFlags().BoolVarP(&indexOptChecksum, "checksum", "C", false, "calculate checksum")
	indexCmd.PersistentFlags().StringVarP(&indexOptMeta, "meta", "m", "", "meta information")
	indexCmd.PersistentFlags().BoolVarP(&indexOptArchive, "archive", "a", false, "index archives")
	indexCmd.PersistentFlags().StringSliceVarP(&indexOptIgnores, "ignore", "i", []string{}, "patterns to ignore")
	indexCmd.PersistentFlags().BoolVarP(&indexOptIndent, "indent", "I", true, "do not indent json")
	indexCmd.PersistentFlags().BoolVarP(&indexOptForce, "force", "f", false, "do not ask user")
	indexCmd.PersistentFlags().BoolVarP(&indexOptNoMIME, "nomime", "M", false, "do not detect mime type")
}

func index(_ *cobra.Command, args []string) error {
	path, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	name := filepath.Base(path)
	if len(args) > 1 {
		name = args[1]
	}
	log.Debugf("indexing \"%s\" as %s", path, name)

	// load the catalog
	t, top, err := loadCatalog(name, path)
	if err != nil {
		log.Fatal(err)
	}

	// build ignore pattern
	var ignPatterns []*regexp.Regexp
	for _, ign := range indexOptIgnores {
		ign = utilities.PatchPattern(ign)
		re, err := regexp.Compile(ign)
		if err != nil {
			log.Fatal(err)
		}
		ignPatterns = append(ignPatterns, re)
	}

	// ensure storage name does not already exist
	for _, storage := range t.Storages {
		if !indexOptForce && name == storage.Name {
			question := fmt.Sprintf("A storage with the name \"%s\" already exists, update it?", name)
			if !utilities.AskUser(question) {
				log.Fatal(fmt.Errorf("user interrupted"))
			}
		}
	}

	// create storage if empty
	if top == nil {
		log.Debugf("creating new storage %s for path %s", name, path)
		// get a new storage
		top = node.NewStorageNode(name, path, filepath.Base(path), indexOptMeta, indexOptTags)
		// and append to tree
		rootTree.Storages = append(rootTree.Storages, top)
	}

	// walk the filesystem
	w := walker.NewWalker(t, indexOptChecksum, indexOptArchive, ignPatterns, indexOptNoMIME)

	t0 := time.Now()
	// spinner
	spinner := pterm.DefaultSpinner.WithRemoveWhenDone(true)
	spinner.Sequence = []string{` ⠋ `, ` ⠙ `, ` ⠹ `, ` ⠸ `, ` ⠼ `, ` ⠴ `, ` ⠦ `, ` ⠧ `, ` ⠇ `, ` ⠏ `}
	spinner.ShowTimer = true
	spinner, err = spinner.Start(fmt.Sprintf("indexing %s", path))
	if err != nil {
		log.Warn(err.Error())
	}

	cnt, size, err := w.Walk(top.ID, path, top, spinner)
	if err == nil {
		log.Debug("stop spinner...")
		err := spinner.Stop()
		if err != nil {
			log.Error(err)
		}
		log.Debug("saving catalog...")
		err = rootCatalog.Save(t)
		if err != nil {
			return err
		}
		hsize := utilities.SizeToHuman(size)
		log.Infof("\"%s\" indexed to \"%s\" (%d entries, %s in %v)", path, rootOptCatalogPath, cnt, hsize, time.Since(t0))
	}
	return err
}

func loadCatalog(storageName string, fsPath string) (*tree.Tree, *node.StorageNode, error) {
	var top *node.StorageNode
	var err error

	if rootTree != nil {
		//log.Debugf("trying to load storage %s", storageName)
		top = rootTree.GetStorageByName(storageName)
		if top != nil {
			log.Debugf("updating storage info for \"%s\"", storageName)
			top.UpdateStorage(fsPath, filepath.Base(fsPath), indexOptMeta, indexOptTags)
		}
	} else {
		// create a new catalog
		rootTree, err = tree.NewTree(version)
		if err != nil {
			return nil, nil, err
		}
	}

	return rootTree, top, nil
}
