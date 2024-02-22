/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/colorme"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
	"path/filepath"
	"strings"
)

const (
	nativeStorageName  = "storage"
	nativeIndentString = "  "
)

// NativeStringer printer struct
type NativeStringer struct {
	theTree *tree.Tree
	mode    *PrintMode
	cm      *colorme.ColorMe
}

func (p *NativeStringer) storageToString(storage *node.StorageNode, pre string) string {
	out := pre
	// "storage"
	out += p.cm.InUnderline(p.cm.InGray(nativeStorageName))
	// the storage name
	out += " "
	out += p.cm.InPurple(fmt.Sprintf("%-20s", storage.GetName()))

	// add attributes
	attrs := storage.GetAttr(p.mode.RawSize, p.mode.Long, p.mode.Extra)
	if len(attrs) > 0 {
		out += " " + AttrsToString(attrs, p.mode, " ")
	}
	return out
}

func (p *NativeStringer) fileToString(n node.Node, pre string) string {
	out := pre

	// name
	name := fmt.Sprintf("%-30s", n.GetName())
	if p.mode.FullPath {
		// full path and storage info
		sto := p.theTree.GetStorageNode(n)
		if sto != nil {
			name = fmt.Sprintf("%-50s", filepath.Join(sto.GetName(), n.GetPath()))
		}
	}
	out += ColorLineByType(name, n, p.mode.InlineColor)

	// add atrributes
	attrs := n.GetAttr(p.mode.RawSize, p.mode.Long, p.mode.Extra)
	if len(attrs) > 0 {
		out += " " + AttrsToString(attrs, p.mode, " ")
	}

	return out
}

// Print prints a node
func (p *NativeStringer) Print(n node.Node, depth int) {
	if n == nil {
		return
	}
	e := p.ToString(n, depth)
	fmt.Println(e.Line)
}

// ToString converts node to string for printing
func (p *NativeStringer) ToString(n node.Node, depth int) *Entry {
	if n == nil {
		return nil
	}
	var entry Entry
	pre := strings.Repeat(nativeIndentString, depth)

	entry.Name = n.GetName()
	entry.Node = n
	if n.GetType() == node.FileTypeStorage {
		entry.Line = p.storageToString(n.(*node.StorageNode), pre)
	} else {
		entry.Line = p.fileToString(n, pre)
	}
	return &entry
}

// PrintPrefix unused
func (p *NativeStringer) PrintPrefix() {}

// PrintSuffix unused
func (p *NativeStringer) PrintSuffix() {}

// NewNativeStringer creates a new native printer
func NewNativeStringer(theTree *tree.Tree, mode *PrintMode) *NativeStringer {
	p := NativeStringer{
		theTree: theTree,
		mode:    mode,
		cm:      colorme.NewColorme(mode.InlineColor),
	}
	return &p
}
