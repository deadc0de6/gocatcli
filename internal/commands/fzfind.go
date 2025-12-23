/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/deadc0de6/gocatcli/internal/log"
	"github.com/deadc0de6/gocatcli/internal/node"
	"github.com/deadc0de6/gocatcli/internal/stringer"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var (
	fzfindCmd = &cobra.Command{
		Use:    "fzfind [<path>]",
		Short:  "Fuzzy find files in the catalog",
		PreRun: preRun(true),
		RunE:   fzFind,
	}

	fzFindOptFormat  string
	fzFindOptDepth   int
	fzFindOptShowAll bool
)

type fzfEntry struct {
	Path    string
	item    node.Node
	storage *node.StorageNode
}

func init() {
	rootCmd.AddCommand(fzfindCmd)

	hlp := fmt.Sprintf("output format (%s)", strings.Join(stringer.GetSupportedFormats(true, false), ","))
	fzfindCmd.PersistentFlags().StringVarP(&fzFindOptFormat, "format", "f", "native", hlp)
	fzfindCmd.PersistentFlags().IntVarP(&fzFindOptDepth, "depth", "D", -1, "max hierarchy depth when printing selected entry")
	fzfindCmd.PersistentFlags().BoolVarP(&fzFindOptShowAll, "all", "a", false, "do not ignore entries starting with a dot")
}

func fzFind(_ *cobra.Command, args []string) error {
	if !formatOk(fzFindOptFormat, true, false) {
		return fmt.Errorf("unsupported format %s", fzFindOptFormat)
	}

	var startPath string
	if len(args) > 0 {
		startPath = args[0]
	}

	// get the stringer
	m := &stringer.PrintMode{
		FullPath:    false,
		Long:        false,
		InlineColor: false,
		RawSize:     false,
		Separator:   separator,
	}
	stringGetter, err := stringer.GetStringer(rootTree, fzFindOptFormat, m)
	if err != nil {
		return err
	}

	// get the base paths for start
	var startNodes []node.Node
	if len(startPath) > 0 {
		startNodes = getStartPaths(startPath)
		if startNodes == nil {
			return fmt.Errorf("no such start path: \"%s\"", startPath)
		}
	} else {
		for _, top := range rootTree.GetStorages() {
			startNodes = append(startNodes, top)
		}
	}

	// list of entries
	var entries []*fzfEntry
	log.Debugf("start nodes: %v", startNodes)
	for _, foundNode := range startNodes {
		entries = append(entries, fzFindFillList(foundNode)...)
	}

	log.Debugf("options contain %d entries", len(entries))

	getItemFunc := func(i int) string {
		return entries[i].Path
	}

	previewFunc := func(i, _, _ int) string {
		if i == -1 {
			return ""
		}
		entry := entries[i]
		var outs []string
		outs = append(outs, fmt.Sprintf("storage: %s", entry.storage.Name))
		outs = append(outs, fmt.Sprintf("path: %s", entry.item.GetPath()))

		entryAttrs := entry.item.GetAttr(m.RawSize, m.Long)
		attrs := stringer.AttrsToString(entryAttrs, m, "\n")

		return strings.Join(outs, "\n") + "\n" + attrs
	}

	// display fzf finder interface
	idx, err := fuzzyfinder.Find(
		entries,
		getItemFunc,
		fuzzyfinder.WithPreviewWindow(previewFunc),
	)
	if err != nil {
		return err
	}

	// print result
	if idx > -1 && idx < len(entries) {
		// list parent directory
		entry := entries[idx]
		log.Debugf("selected entry: %s", entry.Path)

		// get the parent
		hasChildren := false
		typ := entry.item.GetType()
		if typ == node.FileTypeDir || typ == node.FileTypeArchive || typ == node.FileTypeStorage {
			hasChildren = true
		}

		// print the rest
		callback := func(n node.Node, depth int, _ node.Node) bool {
			stringGetter.Print(n, depth+1)
			return true
		}

		stringGetter.PrintPrefix()
		stringGetter.Print(entry.item, 0)
		if hasChildren {
			rootTree.ProcessChildren(entry.item, fzFindOptShowAll, callback, 1)
		}
		stringGetter.PrintSuffix()
	}

	return nil
}

func fzFindFillList(n node.Node) []*fzfEntry {
	var list []*fzfEntry
	top := rootTree.GetStorageNode(n)
	callback := func(n node.Node, _ int, _ node.Node) bool {
		item := &fzfEntry{
			Path:    filepath.Join(top.GetName(), n.GetPath()),
			item:    n,
			storage: top,
		}
		list = append(list, item)
		return true
	}
	rootTree.ProcessChildren(n, fzFindOptShowAll, callback, -1)
	return list
}
