/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"cmp"
	"fmt"
	"gocatcli/internal/node"
	"gocatcli/internal/stringer"
	"slices"

	"github.com/spf13/cobra"
)

var (
	duCmd = &cobra.Command{
		Use:    "du [<path>]",
		Short:  "Disk usage",
		Args:   cobra.MaximumNArgs(1),
		PreRun: preRun(true),
		RunE:   diskUsage,
	}

	duOptRawSize bool
	duOptDepth   int
	duOptSort    bool
)

func init() {
	rootCmd.AddCommand(duCmd)

	duCmd.PersistentFlags().BoolVarP(&duOptRawSize, "raw-size", "S", false, "do not humanize sizes when printing")
	duCmd.PersistentFlags().IntVarP(&duOptDepth, "depth", "D", -1, "max depth")
	duCmd.PersistentFlags().BoolVarP(&duOptSort, "sort", "s", false, "sort by size")
}

func diskUsage(_ *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}

	// get the base paths for start
	var startNodes []node.Node
	if len(path) > 0 {
		startNodes = getStartPaths(path)
		if startNodes == nil {
			return fmt.Errorf("no such start path: \"%s\"", path)
		}
	} else {
		for _, top := range rootTree.GetStorages() {
			startNodes = append(startNodes, top)
		}
	}

	m := &stringer.PrintMode{
		FullPath:    false,
		Long:        true,
		InlineColor: false,
		RawSize:     duOptRawSize,
		Separator:   separator,
	}
	stringer := stringer.NewDuStringer(rootTree, m)
	for _, n := range startNodes {
		var nodes []node.Node
		callback := func(n node.Node, _ int, _ node.Node) bool {
			if !node.IsDir(n) {
				// only dirs
				return true
			}
			if duOptSort {
				nodes = append(nodes, n)
			} else {
				stringer.Print(n, 0, true)
			}
			return true
		}

		rootTree.ProcessChildren(n, true, callback, duOptDepth)
		if len(nodes) > 0 {
			slices.SortFunc(nodes, func(left, right node.Node) int {
				return cmp.Compare(left.GetSize(), right.GetSize())
			})
			for _, n := range nodes {
				stringer.Print(n, 0, true)
			}
		}
		stringer.Print(n, 0, true)
	}

	return nil
}
