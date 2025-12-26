/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package archives

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/deadc0de6/gocatcli/internal/log"

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
	defer func() {
		err := fd.Close()
		if err != nil {
			log.Error(err)
		}
	}()

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
	defer func() {
		err := fd.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	format, stream, err := archiver.Identify(path, fd)
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

	switch archive := format.(type) {
	case archiver.Extractor:
		ctx := context.Background()
		err = archive.Extract(ctx, stream, nil, handler)
		if err != nil {
			return names, err
		}
		log.Debugf("got %d file(s) inside %s", len(names), path)
		return names, nil
	case archiver.Decompressor:
		// no children
		return names, nil
	}

	return nil, fmt.Errorf("cannot read archive content for %s", path)
}
