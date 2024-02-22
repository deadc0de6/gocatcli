/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
	"gocatcli/internal/utils"
	"path/filepath"
)

// DuString printer struct
type DuString struct {
	theTree *tree.Tree
	mode    *PrintMode
}

// Print prints a node
func (p *DuString) Print(n node.Node, depth int, _ bool) {
	if n == nil {
		return
	}
	e := p.ToString(n, depth, true)
	fmt.Println(e.Line)
}

// ToString converts node to string for printing
func (p *DuString) ToString(n node.Node, _ int, _ bool) *Entry {
	if n == nil {
		return nil
	}
	var entry Entry

	prePath := ""
	if !node.IsStorage(n) {
		prePath = p.theTree.GetStorageNode(n).GetName()
	}
	path := filepath.Join(prePath, n.GetPath())
	entry.Name = path
	entry.Node = n
	var size string
	if p.mode.RawSize {
		size = fmt.Sprintf("%d", n.GetSize())
		entry.Line = fmt.Sprintf("%-10s    %s", size, path)
	} else {
		size = utils.SizeToHuman(n.GetSize())
		entry.Line = fmt.Sprintf("%-6s    %s", size, path)
	}
	return &entry
}

// PrintPrefix unused
func (p *DuString) PrintPrefix() {}

// PrintSuffix unused
func (p *DuString) PrintSuffix() {}

// NewDuStringer creates a new native printer
func NewDuStringer(theTree *tree.Tree, mode *PrintMode) *DuString {
	p := DuString{
		theTree: theTree,
		mode:    mode,
	}
	return &p
}
