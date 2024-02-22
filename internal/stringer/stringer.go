/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/node"
	"gocatcli/internal/tree"

	"github.com/TwiN/go-color"
)

const (
	// FormatNative native printing
	FormatNative = "native"
	// FormatCSV csv printing
	FormatCSV = "csv"
	// FormatCSVWithHeader csv with header
	FormatCSVWithHeader = "csv-with-header"
	// FormatScript one-liner script
	FormatScript = "script"
	// FormatTree like unix "tree" command
	FormatTree = "tree"
	// FormatDebug debug output
	FormatDebug = "debug"
)

// Entry entries when traversing the tree
type Entry struct {
	Name string
	Line string
	Node node.Node
}

// PrintMode the type of mode for printing
type PrintMode struct {
	FullPath    bool
	Long        bool
	InlineColor bool
	RawSize     bool
	Separator   string
}

// Stringer interface for printing
type Stringer interface {
	PrintPrefix()
	Print(node.Node, int) // node, depth
	PrintSuffix()
	ToString(node.Node, int) *Entry // node, depth
}

// GetStringer returns a stringer
func GetStringer(tree *tree.Tree, format string, mode *PrintMode) (Stringer, error) {
	var stringGetter Stringer
	switch format {
	case FormatNative:
		stringGetter = NewNativeStringer(tree, mode)
	case FormatCSV:
		stringGetter = NewCSVStringer(tree, mode, false)
	case FormatCSVWithHeader:
		stringGetter = NewCSVStringer(tree, mode, true)
	case FormatScript:
		stringGetter = NewScriptStringer()
	case FormatTree:
		stringGetter = NewTreeStringer(mode)
	case FormatDebug:
		stringGetter = NewDebugStringer(tree, mode)
	default:
		return nil, fmt.Errorf("not such format: %s", format)
	}
	return stringGetter, nil
}

// GetSupportedFormats returns the supported formats
func GetSupportedFormats(treeOk bool, scriptOk bool) []string {
	fmts := []string{
		FormatNative,
		FormatCSV,
		FormatCSVWithHeader,
		FormatDebug,
	}
	if treeOk {
		fmts = append(fmts, FormatTree)
	}
	if scriptOk {
		fmts = append(fmts, FormatScript)
	}
	return fmts
}

// ColorLineByType colors a line by node type
// if inline, inline colors are used
// TODO
func ColorLineByType(line string, n node.Node, inline bool) string {
	var out string
	switch n.GetType() {
	case node.FileTypeDir:
		if inline {
			out = fmt.Sprintf("[blue]%s[-]", line)
		} else {
			out = color.InBlue(line)
		}
	case node.FileTypeArchived:
		if inline {
			out = fmt.Sprintf("[yellow]%s[-]", line)
		} else {
			out = color.InYellow(line)
		}
	case node.FileTypeArchive:
		if inline {
			out = fmt.Sprintf("[red]%s[-]", line)
		} else {
			out = color.InRed(line)
		}
	case node.FileTypeFile:
		fn := n.(*node.FileNode)
		if fn.IsExec() {
			if inline {
				out = fmt.Sprintf("[green]%s[-]", line)
			} else {
				out = color.InGreen(line)
			}
		} else {
			out = line
		}
	default:
		out = line
	}
	return out
}
