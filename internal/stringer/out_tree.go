/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/node"

	"github.com/TwiN/go-color"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// TreeStringer printer struct
type TreeStringer struct {
	mode        *PrintMode
	listOfTrees []*aTree
}

type aTree struct {
	headerLine string
	pterm.LeveledList
}

func (p *TreeStringer) storageToString(storage *node.StorageNode) string {
	out := color.InUnderline(color.InGray(nativeStorageName))
	out += " " + color.InPurple(storage.GetName())
	attrs := storage.GetAttr(false, p.mode.Long, p.mode.Extra)
	if len(attrs) > 0 {
		out += " " + AttrsToString(attrs, p.mode, " ")
	}
	return out
}

func (p *TreeStringer) fileToString(n node.Node) string {
	var out string
	name := fmt.Sprintf("%-20s", n.GetName())
	out += ColorLineByType(name, n, false)
	attrs := n.GetAttr(false, p.mode.Long, p.mode.Extra)
	if len(attrs) > 0 {
		out += " " + AttrsToString(attrs, p.mode, " ")
	}
	return out
}

// Print adds the node to the accumulator
func (p *TreeStringer) Print(n node.Node, depth int) {
	if n == nil {
		// ignore empty node
		return
	}

	isStorage := node.IsStorage(n)
	if isStorage || len(p.listOfTrees) < 1 {
		// new tree
		atree := &aTree{}
		p.listOfTrees = append(p.listOfTrees, atree)
	}

	lastTree := p.listOfTrees[len(p.listOfTrees)-1]
	if isStorage || p.mode.FullPath {
		// add storage or first node as top level
		e := p.ToString(n, depth)
		lastTree.headerLine = e.Line
		return
	}

	// create the item
	e := p.ToString(n, depth)
	item := pterm.LeveledListItem{
		Level: depth,
		Text:  e.Line,
	}

	// append node to the tree
	lastTree.LeveledList = append(lastTree.LeveledList, item)
}

// ToString converts node to string for printing
func (p *TreeStringer) ToString(n node.Node, _ int) *Entry {
	var entry Entry

	if n == nil {
		return nil
	}

	entry.Name = n.GetName()
	entry.Node = n
	if node.IsStorage(n) {
		entry.Line = p.storageToString(n.(*node.StorageNode))
	} else {
		entry.Line = p.fileToString(n)
	}
	return &entry
}

// PrintPrefix unused
func (p *TreeStringer) PrintPrefix() {}

// PrintSuffix print entire tree
func (p *TreeStringer) PrintSuffix() {
	// print each tree
	log.Debugf("number of trees: %d", len(p.listOfTrees))
	for _, atree := range p.listOfTrees {
		if len(atree.headerLine) > 0 {
			fmt.Println(atree.headerLine)
		}
		root := putils.TreeFromLeveledList(atree.LeveledList)
		err := pterm.DefaultTree.WithRoot(root).Render()
		if err != nil {
			log.Error(err)
		}
	}

	// clear
	p.listOfTrees = []*aTree{}
}

// NewTreeStringer creates a new tree printer
func NewTreeStringer(mode *PrintMode) *TreeStringer {
	p := TreeStringer{
		mode: mode,
	}
	return &p
}
