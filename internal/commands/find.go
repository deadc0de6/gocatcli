/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"regexp"
	"strings"
	"time"

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
		// patch pattern
		arg = patchFindPattern(arg)
		// get the pattern to search for
		patt := arg
		re, err := regexp.Compile(patt)
		if err != nil {
			return err
		}
		log.Debugf("search pattern: %s", patt)

		for _, startNode := range startNodes {
			matchNodes(rootTree, startNode, re, stringGetter)
		}
	}

	return nil
}

// find in the tree every node from "startNode" where its name
// matches the pattern "patt"
func matchNodes(t *tree.Tree, startNode node.Node, patt *regexp.Regexp, prt stringer.Stringer) {
	var cnt int64

	t0 := time.Now()
	callback := func(n node.Node, _ int, _ node.Node) bool {
		name := n.GetName()
		log.Debugf("matching name \"%s\" against pattern %v", name, patt)
		ret := patt.MatchString(name)
		if ret {
			log.Debugf("\"%s\" matching \"%v\": %v", name, patt, ret)
			prt.Print(n, 0)
			cnt++
		}
		// always continue
		return true
	}

	prt.PrintPrefix()
	// process all elements of tree
	log.Debugf("processing children and looking for name pattern: %v", patt)
	t.ProcessChildren(startNode, true, callback, -1)
	prt.PrintSuffix()

	log.Debugf("found %d entries matching \"%s\" in %v", cnt, patt.String(), time.Since(t0))
}

// fix pattern
func patchFindPattern(pattern string) string {
	// replace any dot with \.
	patt := strings.ReplaceAll(pattern, ".", "\\.")

	// ensure pattern is enclosed in stars
	if !strings.Contains(patt, "*") {
		ret := fmt.Sprintf(".*%s.*", patt)
		log.Debugf("patched non pattern from \"%s\" to \"%s\"", patt, ret)
		return ret
	}

	// replace all "*" with ".*" for golang pattern
	notDotStar := regexp.MustCompile(`([^\.])\*`)
	ret := notDotStar.ReplaceAllString(patt, "$1.*")

	// replace the first star if any
	if strings.HasPrefix(ret, "*") {
		ret = fmt.Sprintf(".*%s", ret[1:])
	}

	// limit start of line if not star
	if !strings.HasPrefix(ret, ".*") {
		ret = fmt.Sprintf("^%s", ret)
	}

	// limit end of line if not star
	if !strings.HasSuffix(ret, ".*") {
		ret = fmt.Sprintf("%s$", ret)
	}

	log.Debugf("patched pattern from \"%s\" to \"%s\"", pattern, ret)
	return ret
}
