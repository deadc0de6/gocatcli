/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"gocatcli/internal/stringer"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:    "ls [<path>]",
		Short:  "List catalog content",
		Args:   cobra.MaximumNArgs(1),
		PreRun: preRun(true),
		RunE:   list,
	}

	lsOptRecursive bool
	lsOptShowAll   bool
	lsOptFormat    string
	lsOptRawSize   bool
	lsOptLong      bool
	lsOptDepth     int
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolVarP(&lsOptRecursive, "recursive", "r", false, "recursive listing")
	listCmd.PersistentFlags().BoolVarP(&lsOptShowAll, "all", "a", false, "do not ignore entries starting with a dot")
	hlp := fmt.Sprintf("output format (%s)", strings.Join(stringer.GetSupportedFormats(true, true), ","))
	listCmd.PersistentFlags().StringVarP(&lsOptFormat, "format", "f", "native", hlp)
	listCmd.PersistentFlags().BoolVarP(&lsOptRawSize, "raw-size", "S", false, "do not humanize sizes when printing")
	listCmd.PersistentFlags().BoolVarP(&lsOptLong, "long", "l", false, "long listing format")
	listCmd.PersistentFlags().IntVarP(&lsOptDepth, "depth", "D", -1, "max depth")
}

func list(_ *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}
	return ls(path, lsOptFormat, lsOptRecursive, lsOptLong, lsOptRawSize, lsOptShowAll, lsOptDepth)
}

func ls(path string, format string, recursive bool, long bool, rawSize bool, showAll bool, depth int) error {
	if !formatOk(format, true, true) {
		return fmt.Errorf("unsupported format %s", format)
	}
	return listHierarchy(path, recursive, format, showAll, rawSize, long, depth)
}

// list files at "path" or storages if empty
func listHierarchy(path string, recursive bool, format string, showAll bool, rawSize bool, long bool, depth int) error {
	log.Debugf("ls path:%s rec:%v format:%s raw:%v", path, recursive, format, rawSize)
	// get the stringer
	m := &stringer.PrintMode{
		FullPath:    false,
		Long:        long,
		Extra:       long,
		InlineColor: false,
		RawSize:     rawSize,
		Separator:   separator,
	}
	stringGetter, err := stringer.GetStringer(loadedTree, format, m)
	if err != nil {
		return err
	}

	// print prefix
	stringGetter.PrintPrefix()

	log.Debugf("ls path arg: %s (rec:%v)", path, recursive)
	if len(path) < 1 {
		if recursive {
			// list everything recursively
			for _, top := range loadedTree.GetStorages() {
				err := listPrint(stringGetter, top, recursive, showAll, depth)
				if err != nil {
					log.Error(err)
				}
			}
		} else {
			// print the storages only
			// we are intentionally not listing recursively
			// when no storage is selected use find for that
			for _, top := range loadedTree.GetStorages() {
				stringGetter.Print(top, 0)
			}
		}

		// print suffix
		stringGetter.PrintSuffix()
		return nil
	}

	// get the base paths for start
	startNodes := getStartPaths(path)
	if startNodes == nil || len(startNodes) < 1 {
		return fmt.Errorf("no such start path: \"%s\"", path)
	}

	// ls each of the found node
	for _, n := range startNodes {
		log.Debugf("handling found start node \"%s\"", n.GetName())
		err := listPrint(stringGetter, n, recursive, showAll, depth)
		if err != nil {
			log.Error(err)
		}
	}

	// print suffix
	stringGetter.PrintSuffix()

	return nil
}

func listPrint(prt stringer.Stringer, n node.Node, recursive bool, showAll bool, depth int) error {
	// printing the found node
	prt.Print(n, 0)
	if !node.MayHaveChildren(n) {
		return nil
	}

	// print the rest
	callback := func(n node.Node, depth int, _ node.Node) bool {
		prt.Print(n, depth)
		return true
	}

	// for a directory, this will automatically
	// list direct children even if recursive is false
	maxDepth := 0
	if recursive {
		maxDepth = depth
	}
	loadedTree.ProcessChildren(n, showAll, callback, maxDepth)
	return nil
}
