//go:build windows
// +build windows

/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package utilities

import (
	"syscall"
	"unsafe"
)

var (
	lib = "kernel32.dll"
	fn  = "GetDiskFreeSpaceExW"
)

// DiskUsage returns free and total size of disk
func DiskUsage(path string) (uint64, uint64) {
	handle := syscall.MustLoadDLL(lib).MustFindProc(fn)

	var free, total, available int64

	str, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0
	}

	handle.Call(
		uintptr(unsafe.Pointer(str)),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&available)),
	)

	return uint64(free), uint64(total)
}
