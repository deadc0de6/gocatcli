/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"gocatcli/internal/log"
	"gocatcli/internal/navigator"
	"gocatcli/internal/node"
	"gocatcli/internal/stringer"
	"gocatcli/internal/tree"

	"github.com/spf13/cobra"
)

var (
	navCmd = &cobra.Command{
		Use:    "nav [<path>]",
		Short:  "Navigate catalog interactively",
		PreRun: preRun(true),
		RunE:   nav,
	}

	navOptShowAll bool
)

func init() {
	rootCmd.AddCommand(navCmd)

	navCmd.PersistentFlags().BoolVarP(&navOptShowAll, "all", "a", false, "do not ignore entries starting with a dot")
}

func nav(_ *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}

	stringer.DisableColors()

	n := navigator.NewNavigator(getLines(loadedTree), navOptShowAll)

	// get the base paths for start
	startNodes := getStartPaths(path)
	if len(startNodes) > 1 {
		// multiple paths match
		log.Debugf("multiple paths match: \"%s\"", path)
		startNodes = []node.Node{}
	}

	var startPath string
	if len(startNodes) > 0 {
		startPath = startNodes[0].GetPath()
	}
	n.Start(startPath)

	return nil
}

func getLines(t *tree.Tree) func(string, bool) []*stringer.Entry {
	// returns all entries for a specific path (no pattern expected)
	return func(path string, _ bool) []*stringer.Entry {
		var entries []*stringer.Entry

		log.Debugf("nav getting list of files for path \"%s\"", path)

		if len(path) < 1 {
			// return all storages
			log.Debugf("returning all storages...")
			for _, storage := range t.Storages {
				entry := stringer.NewNativeStringer(t, false, false).ToString(storage, 0, false)
				entries = append(entries, entry)
			}
			return entries
		}

		n := loadedTree.GetNodesFromPath(path)
		if len(n) != 1 {
			return nil
		}
		for _, sub := range n[0].GetDirectChildren() {
			printer := stringer.NewNativeStringer(t, false, false)
			entry := printer.ToString(sub, 0, false)
			entries = append(entries, entry)
		}
		return entries
	}
}
