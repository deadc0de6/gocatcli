/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package utilities

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gocatcli/internal/log"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/pterm/pterm"
)

// FileExists returns true if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ChecksumFileContent returns md5 checksum of file content
func ChecksumFileContent(path string) (string, error) {
	if !FileExists(path) {
		return "", fmt.Errorf("%s does not exist", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)
	return hex.EncodeToString(hash[:]), nil
}

// UniqStrings merges and uniq strings in slices
func UniqStrings(slices ...[]string) []string {
	uniq := make(map[string]bool)
	for _, slice := range slices {
		for _, string := range slice {
			uniq[string] = true
		}
	}

	newSlice := make([]string, 0, len(uniq))
	for key := range uniq {
		newSlice = append(newSlice, key)
	}

	return newSlice
}

// SizeToHuman converts size to human readable string
func SizeToHuman(bytes uint64) string {
	str := humanize.Bytes(bytes)
	return strings.ReplaceAll(str, " ", "")
}

// DateToString converts date to string
func DateToString(seconds int64) string {
	dt := time.Unix(seconds, 0)
	return dt.Format("2006-01-02 15:04:05")
}

// HashString creates a 32 bits id from a string
func HashString(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

// HashString64 creates a 64 bits id from a string
func HashString64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// SplitPath splits a path in its components
func SplitPath(path string) []string {
	var paths []string
	if len(path) < 1 {
		return paths
	}
	if path == "/" {
		return []string{path}
	}
	if strings.HasPrefix(path, string(filepath.Separator)) {
		paths = []string{string(filepath.Separator)}
		if len(path) > 1 {
			path = path[1:]
		}
	}

	subs := strings.Split(path, string(filepath.Separator))
	paths = append(paths, subs...)

	log.Debugf("split paths: %#v", paths)
	return paths
}

// AskUser query user for yes/no question
func AskUser(question string) bool {
	resp, _ := pterm.DefaultInteractiveConfirm.Show(question)
	return resp
}

func modeStrToInt(mode []string) int32 {
	var ret int32
	readBit := mode[0]
	if readBit == "r" {
		ret += 4
	}
	writeBit := mode[1]
	if writeBit == "w" {
		ret += 2
	}
	execBit := mode[2]
	if execBit == "x" {
		ret++
	}
	return ret
}

// ModeStrToInt converts mode string to int
func ModeStrToInt(mode string) int32 {
	// -rw-r--r--
	chars := strings.Split(mode, "")
	// drop type indicator
	if len(chars) != 10 {
		log.Warn(fmt.Sprintf("couldn't get mode from %s", mode))
		return 0755
	}
	chars = chars[1:]
	var perm int32
	userVal := modeStrToInt(chars[0:3])
	perm += userVal * 8 * 8
	grpVal := modeStrToInt(chars[3:6])
	perm += grpVal * 8
	othVal := modeStrToInt(chars[6:9])
	perm += othVal
	return perm
}

// NotIn returns true if needle is not in stack
func NotIn(needle string, stack []string) bool {
	for _, element := range stack {
		if needle == element {
			return false
		}
	}
	return true
}

// PatchPattern fix pattern
func PatchPattern(patt string) string {
	// replace any dot with \.
	patt = strings.ReplaceAll(patt, ".", "\\.")

	// ensure pattern is enclosed in stars
	if !strings.Contains(patt, "*") {
		ret := fmt.Sprintf(".*%s.*", patt)
		log.Debugf("patched non pattern from \"%s\" to \"%s\"", patt, ret)
		return ret
	}

	// replace all "*" with ".*" for golang pattern
	notDotStar := regexp.MustCompile(`([^\.])\*`)
	ret := notDotStar.ReplaceAllString(patt, "$1.*")

	// replace the first star if any
	if strings.HasPrefix(ret, "*") {
		ret = fmt.Sprintf(".*%s", ret[1:])
	}

	// limit start of line if not star
	if !strings.HasPrefix(ret, ".*") {
		ret = fmt.Sprintf("^%s", ret)
	}

	// limit end of line if not star
	if !strings.HasSuffix(ret, ".*") {
		ret = fmt.Sprintf("%s$", ret)
	}

	log.Debugf("patched pattern from \"%s\" to \"%s\"", patt, ret)
	return ret
}
