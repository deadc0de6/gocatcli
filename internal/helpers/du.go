//go:build !windows
// +build !windows

/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package helpers

import (
	"syscall"

	"github.com/deadc0de6/gocatcli/internal/log"
)

// DiskUsage returns free and total size of disk
func DiskUsage(path string) (uint64, uint64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.Error(err)
		return 0, 0
	}

	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	return free, total
}
