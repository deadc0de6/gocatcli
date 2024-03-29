/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"gocatcli/internal/catcli"
	"gocatcli/internal/utils"

	"github.com/spf13/cobra"
)

var (
	convertCmd = &cobra.Command{
		Use:    "convert <path>",
		Short:  "Convert catcli catalog",
		Args:   cobra.ExactArgs(1),
		PreRun: preRunDebug,
		RunE:   convert,
	}

	convertOptTo     string
	convertOptIndent bool
)

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringVarP(&convertOptTo, "output", "o", "", "output path")
	convertCmd.PersistentFlags().BoolVarP(&convertOptIndent, "indent", "I", true, "do not indent json")
}

func convert(_ *cobra.Command, args []string) error {
	path := args[0]
	if !utils.FileExists(path) {
		return fmt.Errorf("\"%s\" does not exist", path)
	}

	t, err := catcli.Convert(version, path)
	if err != nil {
		return err
	}

	if len(convertOptTo) > 0 {
		return t.Save(convertOptTo, convertOptIndent)
	}

	content, err := t.Serialize(convertOptIndent)
	if err != nil {
		return err
	}

	fmt.Print(string(content))
	return nil
}
