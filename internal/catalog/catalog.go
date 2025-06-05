package catalog

import (
	"gocatcli/internal/tree"
	"path/filepath"
)

const (
	jsonExt    = ".json"
	tomlExt    = ".toml"
	sqliteExt  = ".sqlite" // TODO
	catalogExt = ".catalog"
)

// Backend catalog backend
type Backend interface {
	Serialize(t *tree.Tree) ([]byte, error)
	Save(path string, t *tree.Tree) error
	LoadTree(path string) (*tree.Tree, error)
}

// Catalog the file catalog
type Catalog struct {
	Path       string
	TheBackend Backend
}

// Serialize the tree
func (c *Catalog) Serialize(t *tree.Tree) ([]byte, error) {
	return c.TheBackend.Serialize(t)
}

// Save tree to file
func (c *Catalog) Save(t *tree.Tree) error {
	return c.TheBackend.Save(c.Path, t)
}

// LoadTree from file
func (c *Catalog) LoadTree() (*tree.Tree, error) {
	return c.TheBackend.LoadTree(c.Path)
}

// NewCatalog creates a new catalog
func NewCatalog(path string) *Catalog {
	var b Backend
	ext := filepath.Ext(path)
	switch ext {
	case catalogExt, jsonExt:
		b = NewJSONBackend()
	case tomlExt:
		b = NewTOMLBackend()
	default:
		// defaults to json
		b = NewJSONBackend()
	}

	c := Catalog{
		Path:       path,
		TheBackend: b,
	}
	return &c
}
