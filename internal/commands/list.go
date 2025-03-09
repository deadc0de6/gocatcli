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
	listCmd.PersistentFlags().IntVarP(&lsOptDepth, "depth", "D", 0, "max depth")
}

func list(_ *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}
	depth := lsOptDepth
	if lsOptRecursive && depth == 0 {
		depth = -1
	}
	return ls(path, lsOptFormat, lsOptLong, lsOptRawSize, lsOptShowAll, depth, false)
}

func ls(path string, format string, long bool, rawSize bool, showAll bool, maxDepth int, printParent bool) error {
	if !formatOk(format, true, true) {
		return fmt.Errorf("unsupported format %s", format)
	}
	return listHierarchy(path, format, showAll, rawSize, long, maxDepth, printParent)
}

// list files at "path" or storages if empty
func listHierarchy(path string, format string, showAll bool, rawSize bool, long bool, maxDepth int, printParent bool) error {
	log.Debugf("ls path:%s format:%s raw:%v maxdepth:%d", path, format, rawSize, maxDepth)
	// get the stringer
	m := &stringer.PrintMode{
		FullPath:    false,
		Long:        long,
		InlineColor: false,
		RawSize:     rawSize,
		Separator:   separator,
	}
	stringGetter, err := stringer.GetStringer(rootTree, format, m)
	if err != nil {
		return err
	}

	// print prefix
	stringGetter.PrintPrefix()

	log.Debugf("ls path arg: %s", path)
	if len(path) < 1 {
		if maxDepth != 0 {
			// list everything recursively
			for _, top := range rootTree.GetStorages() {
				err := listPrint(stringGetter, top, showAll, maxDepth, true)
				if err != nil {
					log.Error(err)
				}
			}
		} else {
			// print the storages only
			// we are intentionally not listing recursively
			// when no storage is selected use find for that
			for _, top := range rootTree.GetStorages() {
				stringGetter.Print(top, 0)
			}
		}

		// print suffix
		stringGetter.PrintSuffix()
		return nil
	}

	// get the base paths for start
	startNodes := getStartPaths(path)
	if len(startNodes) < 1 {
		return fmt.Errorf("no such start path: \"%s\"", path)
	}

	// ls each of the found node
	for _, n := range startNodes {
		log.Debugf("handling found start node \"%s\"", n.GetName())
		err := listPrint(stringGetter, n, showAll, maxDepth, printParent)
		if err != nil {
			log.Error(err)
		}
	}

	// print suffix
	stringGetter.PrintSuffix()

	return nil
}

func listPrint(prt stringer.Stringer, n node.Node, showAll bool, maxDepth int, topToo bool) error {
	var topPrinted bool
	// printing the found node
	if topToo {
		prt.Print(n, 0)
		topPrinted = true
	}

	if !node.MayHaveChildren(n) {
		if !topPrinted {
			prt.Print(n, 0)
		}
		return nil
	}

	// print the rest
	callback := func(n node.Node, depth int, _ node.Node) bool {
		prt.Print(n, depth)
		return true
	}

	// for a directory, this will automatically
	// list direct children even if recursive is false
	rootTree.ProcessChildren(n, showAll, callback, maxDepth)
	return nil
}
