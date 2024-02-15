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
	Name      string      `json:"name"`
	RelPath   string      `json:"relpath"` // to the storage node
	Checksum  string      `json:"md5"`
	Type      FileType    `json:"filetype"`
	Size      uint64      `json:"size"`
	Maccess   int64       `json:"maccess"`
	Children  []*FileNode `json:"children"`
	IndexedAt int64       `json:"ts"`
	StorageID int         `json:"storage_id"`
	Mode      string      `json:"mode"`
	Mime      string      `json:"mime"`
	Extra     string      `json:"extra"` // comma separated list of `<key>:<value>`
	seen      bool        `json:"-"`     // seen tag when updating a storage
}

// StorageNode a storage node
type StorageNode struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Size       uint64      `json:"size"`
	Free       uint64      `json:"free"`
	Total      uint64      `json:"total"`
	IndexedAt  int64       `json:"ts"`
	Type       FileType    `json:"type"`
	Tags       []string    `json:"tags"`
	Meta       string      `json:"meta"`
	TotalFiles uint64      `json:"nb_files"`
	Children   []*FileNode `json:"children"`
}
