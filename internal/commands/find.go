/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/deadc0de6/gocatcli/internal/helpers"
	"github.com/deadc0de6/gocatcli/internal/log"
	"github.com/deadc0de6/gocatcli/internal/node"
	"github.com/deadc0de6/gocatcli/internal/stringer"
	"github.com/deadc0de6/gocatcli/internal/tree"

	"github.com/spf13/cobra"
)

var (
	findCmd = &cobra.Command{
		Use:    "find [<pattern>]",
		Short:  "Find files in the catalog",
		PreRun: preRun(true),
		RunE:   find,
	}

	findOptStart  string
	findOptFormat string
	findOptDepth  int
)

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.PersistentFlags().StringVarP(&findOptStart, "path", "p", "", "start path for find")
	hlp := fmt.Sprintf("output format (%s)", strings.Join(stringer.GetSupportedFormats(false, true), ","))
	findCmd.PersistentFlags().StringVarP(&findOptFormat, "format", "f", "native", hlp)
	findCmd.PersistentFlags().IntVarP(&findOptDepth, "depth", "D", -1, "max depth")
}

func find(_ *cobra.Command, args []string) error {
	// we don't allow "tree" format since
	// multiple node might match and might not be
	// related in the hierarchy, making the output all wrong
	if !formatOk(findOptFormat, false, true) {
		return fmt.Errorf("unsupported format %s", findOptFormat)
	}

	if len(args) < 1 {
		// calling ls when no args are provided
		log.Debugf("running ls recursive...")
		return ls("", findOptFormat, true, false, true, findOptDepth, true)
	}

	// get a stringer to print found nodes
	m := &stringer.PrintMode{
		FullPath:    true,
		Long:        true,
		InlineColor: false,
		RawSize:     false,
		Separator:   separator,
	}
	stringGetter, err := stringer.GetStringer(rootTree, findOptFormat, m)
	if err != nil {
		return err
	}

	// get the base paths for start
	var startNodes []node.Node
	if len(findOptStart) > 0 {
		startNodes = getStartPaths(findOptStart)
		if startNodes == nil {
			return fmt.Errorf("no such start path: \"%s\"", findOptStart)
		}
	} else {
		for _, top := range rootTree.GetStorages() {
			startNodes = append(startNodes, top)
		}
	}

	for _, arg := range args {
		for _, startNode := range startNodes {
			patt := patchFindPattern(arg)
			matchNodes(rootTree, startNode, patt, stringGetter)
		}
	}

	return nil
}

func patchFindPattern(pattern string) string {
	// ensure pattern is enclosed in stars
	if !strings.Contains(pattern, "*") {
		ret := fmt.Sprintf("*%s*", pattern)
		log.Debugf("patched non pattern from \"%s\" to \"%s\"", pattern, ret)
		return ret
	}
	return pattern
}

// find in the tree every node from "startNode" where its name
// matches the pattern "patt"
func matchNodes(t *tree.Tree, startNode node.Node, patt string, prt stringer.Stringer) {
	var cnt int64

	t0 := time.Now()
	callback := func(n node.Node, _ int, _ node.Node) bool {
		name := n.GetName()
		log.Debugf("matching name \"%s\" against pattern \"%s\"", name, patt)
		if helpers.PathMatch(patt, name) {
			log.Debugf("\"%s\" matches \"%s\"", name, patt)
			prt.Print(n, 0)
			cnt++
		}
		// always continue
		return true
	}

	prt.PrintPrefix()
	// process all elements of tree
	log.Debugf("processing children and looking for name pattern: %s", patt)
	t.ProcessChildren(startNode, true, callback, -1)
	prt.PrintSuffix()

	log.Debugf("found %d entries matching \"%s\" in %v", cnt, patt, time.Since(t0))
}
