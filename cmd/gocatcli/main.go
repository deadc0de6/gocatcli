/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package main

import (
	"os"

	"github.com/deadc0de6/gocatcli/internal/commands"
	"github.com/deadc0de6/gocatcli/internal/log"
)

func main() {
	err := commands.Execute()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
