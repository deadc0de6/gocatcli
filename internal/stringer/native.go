/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"
	"path/filepath"
	"strings"

	"github.com/TwiN/go-color"
)

const (
	nativeStorageName  = "storage"
	nativeIndentString = "  "
)

// NativeStringer printer struct
type NativeStringer struct {
	theTree  *tree.Tree
	rawSize  bool
	fullInfo bool
}

func (p *NativeStringer) storageToString(storage *node.StorageNode, pre string) string {
	out := pre
	out += color.InUnderline(color.InGray(nativeStorageName))
	out += " " + color.InPurple(storage.GetName())
	attrs := storage.GetAttr(p.rawSize, p.fullInfo)
	if len(attrs) > 0 {
		out += " " + AttrsToString(true, attrs, " ")
	}
	return out
}

func (p *NativeStringer) fileToString(n node.Node, pre string, fullPath bool) string {
	out := pre

	name := fmt.Sprintf("%-5s", n.GetName())
	if p.fullInfo {
		// full path and storage info
		sto := p.theTree.GetStorageNode(n)
		if fullPath && sto != nil {
			name = fmt.Sprintf("%-20s", filepath.Join(sto.GetName(), n.GetPath()))
		}
	}

	out += ColorByType(name, n, false)

	// add atrributes
	attrs := n.GetAttr(p.rawSize, p.fullInfo)
	if len(attrs) > 0 {
		out += " " + AttrsToString(p.fullInfo, attrs, " ")
	}

	return out
}

// Print prints a node
func (p *NativeStringer) Print(n node.Node, depth int, fullPath bool) {
	if n == nil {
		return
	}
	e := p.ToString(n, depth, fullPath)
	fmt.Println(e.Line)
}

// ToString converts node to string for printing
func (p *NativeStringer) ToString(n node.Node, depth int, fullPath bool) *Entry {
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
		entry.Line = p.fileToString(n, pre, fullPath)
	}
	return &entry
}

// PrintPrefix unused
func (p *NativeStringer) PrintPrefix() {}

// PrintSuffix unused
func (p *NativeStringer) PrintSuffix() {}

// NewNativeStringer creates a new native printer
func NewNativeStringer(theTree *tree.Tree, rawSize bool, fullInfo bool) *NativeStringer {
	p := NativeStringer{
		theTree:  theTree,
		rawSize:  rawSize,
		fullInfo: fullInfo,
	}
	return &p
}
