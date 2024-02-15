package walker

import (
	"gocatcli/internal/log"
	"os"

	"github.com/h2non/filetype"
)

var (
	headerSize = 512
)

func getMime(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Debugf("filetype open error: %v", err)
		return ""
	}
	head := make([]byte, headerSize)
	_, err = file.Read(head)
	if err != nil {
		log.Debugf("filetype read error: %v", err)
		return ""
	}

	log.Debugf("getting mime of %s", path)
	m, err := filetype.Match(head)
	if err != nil {
		log.Debugf("filetype match error: %v", err)
		return ""
	}
	log.Debugf("\"%s\" mime: %s", path, m.MIME.Value)
	return m.MIME.Value
}
