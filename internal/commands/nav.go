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
	navOptLong    bool
)

func init() {
	rootCmd.AddCommand(navCmd)

	navCmd.PersistentFlags().BoolVarP(&navOptShowAll, "all", "a", false, "do not ignore entries starting with a dot")
	navCmd.PersistentFlags().BoolVarP(&navOptLong, "long", "l", false, "long listing format")
}

func nav(_ *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}

	n := navigator.NewNavigator(callback(loadedTree), navOptShowAll, navOptLong)

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

func callback(t *tree.Tree) func(string, bool, bool) (bool, []*stringer.Entry) {
	// returns all entries for a specific path (no pattern expected)
	return func(path string, showHidden bool, longMode bool) (bool, []*stringer.Entry) {
		var entries []*stringer.Entry

		log.Debugf("nav getting list of files for path \"%s\"", path)
		log.Debugf("showhiden: %v", showHidden)
		log.Debugf("long mode: %v", longMode)

		m := &stringer.PrintMode{
			FullPath:    true,
			Long:        longMode,
			Extra:       longMode,
			InlineColor: true,
			RawSize:     false,
			Separator:   separator,
		}
		printer := stringer.NewNativeStringer(t, m)
		if len(path) < 1 {
			// return all storages
			log.Debugf("returning all storages...")
			for _, storage := range t.Storages {
				entry := printer.ToString(storage, 0)
				entries = append(entries, entry)
			}
			return true, entries
		}

		founds := loadedTree.GetNodesFromPath(path)
		if len(founds) != 1 {
			// nothing there
			return false, nil
		}

		// fill the list
		callback := func(n node.Node, depth int, _ node.Node) bool {
			sub := printer.ToString(n, 0)
			entries = append(entries, sub)
			return true
		}
		loadedTree.ProcessChildren(founds[0], showHidden, callback, 0)
		return false, entries
	}
}
