/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package stringer

import (
	"fmt"
	"gocatcli/internal/colorme"
	"gocatcli/internal/node"
	"gocatcli/internal/utils"
	"sort"
	"strings"
)

var (
	topAttrs     = []string{"mode", "type", "size", "maccess"}
	extraAttrs   = []string{"indexed", "children", "checksum"}
	childrenAttr = "children"
)

func attrToStringColored(key string, value string, cm *colorme.ColorMe) string {
	var line string
	if key == "date" {
		line = cm.InBlue(value)
	} else if key == "maccess" {
		line = cm.InBlue(value)
	} else if key == "mode" {
		line = cm.InYellow(value)
	} else if key == "size" {
		line = cm.InGreen(fmt.Sprintf("%6s", value))
	} else if key == "type" {
		line = cm.InRed(fmt.Sprintf("%-4s", value))
	} else {
		line = cm.InGray(value)
	}

	return line
}

func getAttr(attrs map[string]string, key string) string {
	val, ok := attrs[key]
	if !ok {
		return ""
	}

	return val
}

func getMoreAttrs(attrs map[string]string, notThose []string, cm *colorme.ColorMe) []string {
	var outs []string

	skipChildren := false
	if getAttr(attrs, "type") == node.FileTypeFile || getAttr(attrs, "type") == node.FileTypeArchive || getAttr(attrs, "type") == node.FileTypeArchived {
		skipChildren = true
	}

	// get the extra first
	for _, key := range extraAttrs {
		if key == childrenAttr && skipChildren {
			continue
		}
		val := getAttr(attrs, key)
		if len(val) < 1 {
			continue
		}
		line := fmt.Sprintf("%s:%s", cm.InGray(key), attrs[key])
		outs = append(outs, line)
		notThose = append(notThose, key)
	}

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		if utils.NotIn(k, notThose) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, key := range keys {
		if key == childrenAttr && skipChildren {
			continue
		}
		val := getAttr(attrs, key)
		if len(val) < 1 {
			continue
		}
		line := fmt.Sprintf("%s:%s", cm.InGray(key), attrs[key])
		outs = append(outs, line)
	}
	return outs
}

// AttrsToString converts attributes to string
func AttrsToString(attrs map[string]string, mode *PrintMode, joiner string) string {
	var outs []string

	if !mode.Long {
		return strings.Join(outs, joiner)
	}

	cm := colorme.NewColorme(mode.InlineColor)
	for _, key := range topAttrs {
		val := getAttr(attrs, key)
		if len(val) > 0 {
			outs = append(outs, attrToStringColored(key, val, cm))
		}
	}
	outs = append(outs, getMoreAttrs(attrs, topAttrs, cm)...)

	return strings.Join(outs, joiner)
}
