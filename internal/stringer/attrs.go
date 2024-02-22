/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/node"
	"gocatcli/internal/utils"
	"sort"
	"strings"

	"github.com/TwiN/go-color"
)

// TODO add colors

var (
	topAttrs     = []string{"mode", "type", "size", "maccess"}
	childrenAttr = "children"
)

func getAttr(attrs map[string]string, key string) string {
	val, ok := attrs[key]
	if !ok {
		return ""
	}

	return val
}

func getMoreAttrs(attrs map[string]string, notThose []string) []string {
	var outs []string

	skipChildren := false
	if getAttr(attrs, "type") == node.FileTypeFile || getAttr(attrs, "type") == node.FileTypeArchive || getAttr(attrs, "type") == node.FileTypeArchived {
		skipChildren = true
	}

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if key == childrenAttr && skipChildren {
			continue
		}
		if utils.NotIn(key, notThose) {
			val := getAttr(attrs, key)
			if len(val) > 0 {
				outs = append(outs, fmt.Sprintf("%s:%s", key, attrs[key]))
			}
		}
	}
	return outs
}

// AttrsToString converts attributes to string
func AttrsToString(attrs map[string]string, mode *PrintMode, joiner string) string {
	var outs []string

	if !mode.Long {
		return strings.Join(outs, joiner)
	}

	for _, attr := range topAttrs {
		val := getAttr(attrs, attr)
		if len(val) > 0 {
			outs = append(outs, color.InGray(val))
		}
	}
	outs = append(outs, getMoreAttrs(attrs, topAttrs)...)

	if !mode.Extra {
		return strings.Join(outs, joiner)
	}

	for _, attr := range topAttrs {
		val := getAttr(attrs, attr)
		if len(val) > 0 {
			outs = append(outs, color.InGray(val))
		}
	}

	return strings.Join(outs, joiner)
}
