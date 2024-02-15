/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
)

// DebugStringer printer struct
type DebugStringer struct {
	theTree *tree.Tree
	rawSize bool
}

// Print prints a node
func (p *DebugStringer) Print(n node.Node, depth int, fullPath bool) {
	if n == nil {
		return
	}
	e := p.ToString(n, depth, fullPath)
	fmt.Println(e.Line)
}

// ToString converts node to string for printing
func (p *DebugStringer) ToString(n node.Node, _ int, _ bool) *Entry {
	if n == nil {
		return nil
	}
	var entry Entry
	entry.Name = n.GetName()
	entry.Node = n
	entry.Line = fmt.Sprintf("%#v", n)
	return &entry
}

// PrintPrefix unused
func (p *DebugStringer) PrintPrefix() {}

// PrintSuffix unused
func (p *DebugStringer) PrintSuffix() {}

// NewDebugStringer creates a new debug printer
func NewDebugStringer(theTree *tree.Tree, rawSize bool) *DebugStringer {
	p := DebugStringer{
		theTree: theTree,
		rawSize: rawSize,
	}
	return &p
}
