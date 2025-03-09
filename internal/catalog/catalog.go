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

type Backend interface {
	Serialize(t *tree.Tree) ([]byte, error)
	Save(path string, t *tree.Tree) error
	LoadTree(path string) (*tree.Tree, error)
}

type Catalog struct {
	Path string
	B    Backend
}

func (c *Catalog) Serialize(t *tree.Tree) ([]byte, error) {
	return c.B.Serialize(t)
}

func (c *Catalog) Save(t *tree.Tree) error {
	return c.B.Save(c.Path, t)
}

func (c *Catalog) LoadTree() (*tree.Tree, error) {
	return c.B.LoadTree(c.Path)
}

func NewCatalog(path string) *Catalog {
	var b Backend
	ext := filepath.Ext(path)
	if ext == catalogExt {
		b = NewTOMLBackend()
	} else if ext == jsonExt {
		b = NewJSONBackend()
	} else if ext == tomlExt {
		b = NewTOMLBackend()
	}
	c := Catalog{
		Path: path,
		B:    b,
	}
	return &c
}
