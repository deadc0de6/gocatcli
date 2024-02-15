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

// PTreeStringer printer struct
type PTreeStringer struct {
	long        bool
	listOfTrees []*aTree
}

type aTree struct {
	headerLine string
	pterm.LeveledList
}

func (p *PTreeStringer) storageToString(storage *node.StorageNode) string {
	out := color.InUnderline(color.InGray(nativeStorageName))
	out += " " + color.InPurple(storage.GetName())
	attrs := storage.GetAttr(false, p.long)
	if len(attrs) > 0 {
		out += " " + AttrsToString(true, attrs, " ")
	}
	return out
}

func (p *PTreeStringer) fileToString(n node.Node, _ bool) string {
	var out string
	name := fmt.Sprintf("%-20s", n.GetName())
	out += ColorByType(name, n, false)
	attrs := n.GetAttr(false, p.long)
	if len(attrs) > 0 {
		out += " " + AttrsToString(p.long, attrs, " ")
	}
	return out
}

// Print adds the node to the accumulator
func (p *PTreeStringer) Print(n node.Node, depth int, fullPath bool) {
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
	if isStorage || fullPath {
		// add storage or first node as top level
		e := p.ToString(n, depth, fullPath)
		lastTree.headerLine = e.Line
		return
	}

	// create the item
	e := p.ToString(n, depth, fullPath)
	item := pterm.LeveledListItem{
		Level: depth,
		Text:  e.Line,
	}

	// append node to the tree
	lastTree.LeveledList = append(lastTree.LeveledList, item)
}

// ToString converts node to string for printing
func (p *PTreeStringer) ToString(n node.Node, _ int, fullPath bool) *Entry {
	var entry Entry

	if n == nil {
		return nil
	}

	entry.Name = n.GetName()
	entry.Node = n
	if node.IsStorage(n) {
		entry.Line = p.storageToString(n.(*node.StorageNode))
	} else {
		entry.Line = p.fileToString(n, fullPath)
	}
	return &entry
}

// PrintPrefix unused
func (p *PTreeStringer) PrintPrefix() {}

// PrintSuffix print entire tree
func (p *PTreeStringer) PrintSuffix() {
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

// NewPTreeStringer creates a new ptree printer
func NewPTreeStringer(long bool) *PTreeStringer {
	p := PTreeStringer{
		long: long,
	}
	return &p
}
