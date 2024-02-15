/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package archives

import (
	"context"
	"gocatcli/internal/log"
	"io/fs"
	"os"

	"github.com/mholt/archiver/v4"
)

// ArchivedFile file inside an archive
type ArchivedFile struct {
	FileInfo fs.FileInfo
	Path     string
}

// IsArchive returns true if file pointed by path is a supported archive
func IsArchive(path string) bool {
	fd, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return false
	}
	defer fd.Close()

	_, _, err = archiver.Identify(path, fd)
	return err != archiver.ErrNoMatch
}

// GetFiles return the list of files in archive
func GetFiles(path string) ([]*ArchivedFile, error) {
	var names []*ArchivedFile

	fd, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return names, err
	}
	defer fd.Close()

	format, _, err := archiver.Identify(path, fd)
	if err == archiver.ErrNoMatch {
		log.Debugf("file \"%s\" is not an archive", path)
		return names, nil
	}
	if err != nil {
		return names, err
	}
	log.Debugf("process archive \"%s\" as \"%s\"", path, format.Name())

	handler := func(_ context.Context, f archiver.File) error {
		arc := ArchivedFile{
			FileInfo: f.FileInfo,
			Path:     f.NameInArchive,
		}
		names = append(names, &arc)
		return nil
	}

	ext := format.(archiver.Extractor)
	ctx := context.Background()
	err = ext.Extract(ctx, fd, nil, handler)
	if err != nil {
		return names, err
	}

	log.Debugf("got %d file(s) inside %s", len(names), path)
	return names, nil
}
