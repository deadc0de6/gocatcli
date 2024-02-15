/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package catcli

import "encoding/json"

var (
	//nodeTypeTop     = "top"
	//nodeTypeStorage = "storage"
	//nodeTypeMeta    = "meta"
	nodeTypeFile    = "file"
	nodeTypeDir     = "dir"
	nodeTypeArchive = "arc"
)

// Top top node
type Top struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Children []json.RawMessage `json:"children"`
}

// Meta meta node for storage
type Meta struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	Meta MetaAttr `json:"attr"`
}

// MetaAttr meta node attribute
type MetaAttr struct {
	Access         int64  `json:"access"`
	AccessVersion  string `json:"access_version"`
	Created        int64  `json:"created"`
	CreatedVersion string `json:"created_version"`
}

// Storage storage node
type Storage struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Free      int64   `json:"free"`
	Size      int64   `json:"size"`
	Total     int64   `json:"total"`
	TimeStamp int     `json:"ts"`
	Attr      string  `json:"attr"` // comma separated list of tags (called meta in catcli)
	Children  []*Node `json:"children"`
}

// Node file|dir|arc
type Node struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	MAccess  float32 `json:"maccess"`
	MD5      string  `json:"md5"`
	RelPath  string  `json:"relpath"`
	Size     int     `json:"size"`
	Archive  string  `json:"archive,omitempty"`
	Children []*Node `json:"children"`
}
