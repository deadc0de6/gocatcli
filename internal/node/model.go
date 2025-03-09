/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package node

const (
	// FileTypeFile a file
	FileTypeFile = "file"
	// FileTypeDir a directory
	FileTypeDir = "dir"
	// FileTypeStorage a storage
	FileTypeStorage = "storage"
	// FileTypeArchive a file archive
	FileTypeArchive = "archive"
	// FileTypeArchived an archived file
	FileTypeArchived = "archived"
)

// FileNode a file node
type FileNode struct {
	ID        string      `json:"id" toml:"id"`
	Name      string      `json:"name" toml:"name"`
	RelPath   string      `json:"relpath" toml:"relpath"` // to the storage node
	Checksum  string      `json:"md5" toml:"md5"`
	Type      FileType    `json:"filetype" toml:"filetype"`
	Size      uint64      `json:"size" toml:"size"`
	Maccess   int64       `json:"maccess" toml:"maccess"`
	Children  []*FileNode `json:"children" toml:"children"`
	IndexedAt int64       `json:"ts" toml:"ts"`
	StorageID int         `json:"storage_id" toml:"storage_id"`
	Mode      string      `json:"mode" toml:"mode"`
	Mime      string      `json:"mime" toml:"mime"`
	Extra     string      `json:"extra" toml:"extra"` // comma separated list of `<key>:<value>`
	seen      bool        `json:"-" toml:"-"`         // seen tag when updating a storage
}

// StorageNode a storage node
type StorageNode struct {
	ID         int         `json:"id" toml:"id"`
	Name       string      `json:"name" toml:"name"`
	Path       string      `json:"path" toml:"path"`
	Size       uint64      `json:"size" toml:"size"`
	Free       uint64      `json:"free" toml:"free"`
	Total      uint64      `json:"total" toml:"total"`
	IndexedAt  int64       `json:"ts" toml:"ts"`
	Type       FileType    `json:"type" toml:"type"`
	Tags       []string    `json:"tags" toml:"tags"`
	Meta       string      `json:"meta" toml:"meta"`
	TotalFiles uint64      `json:"nb_files" toml:"nb_files"`
	Children   []*FileNode `json:"children" toml:"children"`
}
