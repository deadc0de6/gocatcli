/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"github.com/deadc0de6/gocatcli/internal/stringer"

	"github.com/spf13/cobra"
)

var (
	treeCmd = &cobra.Command{
		Use:    "tree [<path>]",
		Short:  "List catalog content as a tree",
		PreRun: preRun(true),
		RunE:   treeView,
	}

	treeOptShowAll bool
	treeOptRawSize bool
	treeOptLong    bool
	treeOptDepth   int
)

func init() {
	rootCmd.AddCommand(treeCmd)

	treeCmd.PersistentFlags().BoolVarP(&treeOptShowAll, "all", "a", false, "do not ignore entries starting with a dot")
	treeCmd.PersistentFlags().BoolVarP(&treeOptRawSize, "raw-size", "S", false, "do not humanize sizes when printing")
	treeCmd.PersistentFlags().BoolVarP(&treeOptLong, "long", "l", false, "long listing format")
	treeCmd.PersistentFlags().IntVarP(&treeOptDepth, "depth", "D", -1, "max depth")
}

func treeView(_ *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}
	return ls(path, stringer.FormatTree, treeOptLong, treeOptRawSize, treeOptShowAll, treeOptDepth, true)
}
