/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"

	"github.com/deadc0de6/gocatcli/internal/node"
	"github.com/deadc0de6/gocatcli/internal/tree"
)

// DebugStringer printer struct
type DebugStringer struct {
	theTree *tree.Tree
	mode    *PrintMode
}

// Print prints a node
func (p *DebugStringer) Print(n node.Node, depth int) {
	if n == nil {
		return
	}
	e := p.ToString(n, depth)
	fmt.Println(e.Line)
}

// ToString converts node to string for printing
func (p *DebugStringer) ToString(n node.Node, _ int) *Entry {
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
func NewDebugStringer(theTree *tree.Tree, mode *PrintMode) *DebugStringer {
	p := DebugStringer{
		theTree: theTree,
		mode:    mode,
	}
	return &p
}
