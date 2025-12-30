/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"

	"github.com/deadc0de6/gocatcli/internal/colorme"
	"github.com/deadc0de6/gocatcli/internal/node"
	"github.com/deadc0de6/gocatcli/internal/tree"
)

const (
	filenameStorageName = "storage"
)

// FileNameStringer printer struct
type FileNameStringer struct {
	theTree *tree.Tree
	mode    *PrintMode
	cm      *colorme.ColorMe
}

func (p *FileNameStringer) storageToString(storage *node.StorageNode) string {
	// "storage"
	out := p.cm.InUnderline(p.cm.InGray(filenameStorageName))
	// the storage name
	out += " "
	out += p.cm.InPurple(storage.GetName())
	return out
}

func (p *FileNameStringer) fileToString(n node.Node) string {
	return ColorLineByType(n.GetName(), n, p.mode.InlineColor)
}

// Print prints a node
func (p *FileNameStringer) Print(n node.Node, depth int) {
	if n == nil {
		return
	}
	e := p.ToString(n, depth)
	fmt.Println(e.Line)
}

// ToString converts node to string for printing
func (p *FileNameStringer) ToString(n node.Node, _ int) *Entry {
	if n == nil {
		return nil
	}
	var entry Entry

	entry.Name = n.GetName()
	entry.Node = n
	if n.GetType() == node.FileTypeStorage {
		entry.Line = p.storageToString(n.(*node.StorageNode))
	} else {
		entry.Line = p.fileToString(n)
	}
	return &entry
}

// PrintPrefix unused
func (p *FileNameStringer) PrintPrefix() {}

// PrintSuffix unused
func (p *FileNameStringer) PrintSuffix() {}

// NewFileNameStringer creates a new native printer
func NewFileNameStringer(theTree *tree.Tree, mode *PrintMode) *FileNameStringer {
	p := FileNameStringer{
		theTree: theTree,
		mode:    mode,
		cm:      colorme.NewColorme(mode.InlineColor),
	}
	return &p
}
