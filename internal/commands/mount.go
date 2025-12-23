/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"github.com/deadc0de6/gocatcli/internal/fuser"

	"github.com/spf13/cobra"
)

var (
	mountCmd = &cobra.Command{
		Use:    "mount [<path>]",
		Short:  "Mount catalog using fuse",
		Args:   cobra.ExactArgs(1),
		PreRun: preRun(true),
		RunE:   mount,
	}
)

func init() {
	rootCmd.AddCommand(mountCmd)
}

func mount(_ *cobra.Command, args []string) error {
	path := args[0]
	err := fuser.Mount(rootTree, path, rootOptDebugMode)
	return err
}
