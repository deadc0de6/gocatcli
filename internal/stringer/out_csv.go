/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"strings"

	"github.com/deadc0de6/gocatcli/internal/node"
	"github.com/deadc0de6/gocatcli/internal/tree"
	"github.com/deadc0de6/gocatcli/internal/utilities"
)

var (
	header = []string{
		"name",
		"type",
		"path",
		"size",
		"indexed_at",
		"maccess",
		"checksum",
		"nbfiles",
		"free_space",
		"total_space",
		"meta",
		"storage",
	}
)

// CSVStringer the CSV stringer
type CSVStringer struct {
	theTree    *tree.Tree
	mode       *PrintMode
	withHeader bool
}

func (p *CSVStringer) getSize(sz uint64) string {
	if p.mode.RawSize {
		return fmt.Sprintf("%d", sz)
	}
	return utilities.SizeToHuman(sz)
}

func (p *CSVStringer) storageToString(storage *node.StorageNode) string {
	var fields []string
	fields = append(fields, storage.Name)
	fields = append(fields, string(storage.Type))
	fields = append(fields, storage.GetPath())
	fields = append(fields, p.getSize(storage.Size))
	fields = append(fields, utilities.DateToString(storage.IndexedAt))
	fields = append(fields, "") // maccess
	fields = append(fields, "") // checksum
	fields = append(fields, fmt.Sprintf("%d", storage.TotalFiles))
	fields = append(fields, p.getSize(storage.Free))
	fields = append(fields, p.getSize(storage.Total))
	fields = append(fields, storage.Meta)
	fields = append(fields, storage.Name) // storage (self)
	return strings.Join(fields, p.mode.Separator)
}

func (p *CSVStringer) fileToString(n *node.FileNode) string {
	var fields []string
	fields = append(fields, n.Name)
	fields = append(fields, string(n.Type))
	fields = append(fields, n.GetPath())
	fields = append(fields, p.getSize(n.Size))
	fields = append(fields, utilities.DateToString(n.IndexedAt))
	maccess := utilities.DateToString(n.Maccess)
	fields = append(fields, maccess)
	fields = append(fields, string(n.Checksum))
	fields = append(fields, fmt.Sprintf("%d", len(n.Children)))
	fields = append(fields, "") // free_space
	fields = append(fields, "") // total_space
	fields = append(fields, "") // meta
	sto := p.theTree.GetStorageNode(n)
	if sto != nil {
		fields = append(fields, sto.GetName()) // storage
	} else {
		fields = append(fields, "")
	}

	return strings.Join(fields, p.mode.Separator)
}

// ToString converts node to csv for printing
func (p *CSVStringer) ToString(n node.Node, _ int) *Entry {
	var entry Entry

	entry.Name = n.GetName()
	entry.Node = n
	if n.GetType() == node.FileTypeStorage {
		entry.Line = p.storageToString(n.(*node.StorageNode))
	} else {
		entry.Line = p.fileToString(n.(*node.FileNode))
	}
	return &entry
}

// PrintPrefix prints the header
func (p *CSVStringer) PrintPrefix() {
	if !p.withHeader {
		return
	}
	fmt.Println(strings.Join(header, p.mode.Separator))
}

// PrintSuffix unused
func (p *CSVStringer) PrintSuffix() {}

// Print prints a node
func (p *CSVStringer) Print(n node.Node, depth int) {
	e := p.ToString(n, depth)
	fmt.Println(e.Line)
}

// NewCSVStringer creates a new CSV printer
func NewCSVStringer(t *tree.Tree, mode *PrintMode, withHeader bool) *CSVStringer {
	p := CSVStringer{
		theTree:    t,
		mode:       mode,
		withHeader: withHeader,
	}
	return &p
}
