package catalog

import (
	"encoding/json"
	"gocatcli/internal/log"
	"gocatcli/internal/tree"
	"os"
	"time"
)

var (
	jsonIndent = true
)

// JSONBackend the JSON backend
type JSONBackend struct{}

// Serialize gets the tree as string
func (b *JSONBackend) Serialize(t *tree.Tree) ([]byte, error) {
	t.Updated = time.Now().Unix()

	var content []byte
	var err error
	if jsonIndent {
		content, err = json.MarshalIndent(t, "", "  ")
	} else {
		content, err = json.Marshal(t)
	}
	if err != nil {
		log.Debugf("marshal failed: %v", err)
		return nil, err
	}

	return content, err
}

// Save saves a tree to json
func (b *JSONBackend) Save(path string, t *tree.Tree) error {
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

// LoadTree loads a tree from json
func (b *JSONBackend) LoadTree(path string) (*tree.Tree, error) {
	log.Debugf("loading catalog from %s", path)
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var tree tree.Tree
	err = json.NewDecoder(fd).Decode(&tree)
	if err != nil {
		return nil, err
	}

	return &tree, nil
}

// NewJSONBackend creates a new json backend
func NewJSONBackend() *JSONBackend {
	b := &JSONBackend{}
	return b
}
