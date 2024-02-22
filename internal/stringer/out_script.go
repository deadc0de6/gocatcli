/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/node"
)

const (
	scriptPrefix = "op=file; source=/media/mnt; ${op} "
)

// ScriptStringer printer struct
type ScriptStringer struct{}

// Print prints a node
func (p *ScriptStringer) Print(n node.Node, _ int) {
	if n == nil {
		return
	}
	typ := n.GetType()
	if typ == node.FileTypeArchived {
		return
	}
	if typ == node.FileTypeStorage {
		return
	}
	fmt.Printf("\"${source}/%s\" ", n.GetPath())
}

// ToString unsupported
func (p *ScriptStringer) ToString(node.Node, int) *Entry {
	return nil
}

// PrintPrefix unused
func (p *ScriptStringer) PrintPrefix() {
	fmt.Print(scriptPrefix)
}

// PrintSuffix unused
func (p *ScriptStringer) PrintSuffix() {
	fmt.Println()
}

// NewScriptStringer creates a new script printer
func NewScriptStringer() *ScriptStringer {
	return &ScriptStringer{}
}
