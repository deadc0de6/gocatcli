package catalog

import (
	"bytes"
	"gocatcli/internal/log"
	"gocatcli/internal/tree"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// TOMLBackend the toml backend
type TOMLBackend struct{}

// Serialize gets the tree as string
func (b *TOMLBackend) Serialize(t *tree.Tree) ([]byte, error) {
	t.Updated = time.Now().Unix()

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(t)
	if err != nil {
		log.Debugf("marshal failed: %v", err)
		return nil, err
	}
	return buf.Bytes(), err
}

// Save saves a tree to toml
func (b *TOMLBackend) Save(path string, t *tree.Tree) error {
	log.Debug("serialize tree...")
	content, err := b.Serialize(t)
	if err != nil {
		return err
	}

	log.Debugf("write tree to \"%s\"...", path)
	err = os.WriteFile(path, content, os.ModePerm)
	if err != nil {
		return err
	}
	log.Debugf("tree saved to \"%s\"", path)
	return nil
}

// LoadTree loads a tree from toml
func (b *TOMLBackend) LoadTree(path string) (*tree.Tree, error) {
	log.Debugf("loading catalog from %s", path)
	var tree tree.Tree
	_, err := toml.DecodeFile(path, &tree)
	if err != nil {
		return nil, err
	}

	return &tree, nil
}

// NewTOMLBackend creates a new toml backend
func NewTOMLBackend() *TOMLBackend {
	b := &TOMLBackend{}
	return b
}
