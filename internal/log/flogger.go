/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package log

import (
	"os"

	"github.com/caarlos0/log"
)

// appendToFile appends content (string) to file
func appendToFile(path string, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	_, err = f.WriteString(content)
	return err
}

// ToFile saves log to file
func ToFile(path string, content string) {
	err := appendToFile(path, content+"\n")
	if err != nil {
		Error(err)
	}
}
