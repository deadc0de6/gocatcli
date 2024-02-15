/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"gocatcli/internal/log"
	"gocatcli/internal/stringer"
	"gocatcli/internal/tree"
	"gocatcli/internal/utils"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version    = "0.1"
	myName     = "gocatcli"
	defCatalog = "gocatcli.catalog"
	loadedTree *tree.Tree
	separator  = ","

	rootCmd = &cobra.Command{
		Use:     "gocatcli",
		Short:   "gocatcli - filesystem indexer",
		Long:    `The command line catalog tool for your offline data`,
		Version: version,
	}

	rootOptCatalogPath string
	rootOptDebugMode   bool
)

func init() {
	// env variables
	viper.SetEnvPrefix(strings.ToUpper(myName))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// flags
	defCatalogPath := viper.GetString("CATALOG")
	if len(defCatalogPath) < 1 {
		defCatalogPath = defCatalog
	}
	rootCmd.PersistentFlags().StringVarP(&rootOptCatalogPath, "catalog", "c", defCatalogPath, "catalog file path")
	rootCmd.PersistentFlags().BoolVarP(&rootOptDebugMode, "debug", "d", viper.GetBool("DEBUG"), "enable debug mode")
}

func preRunDebug(*cobra.Command, []string) {
	if rootOptDebugMode {
		log.DebugMode = true
	}
}

func preRun(loadCatalogFatal bool) func(*cobra.Command, []string) {
	return func(ccmd *cobra.Command, args []string) {
		var err error

		preRunDebug(ccmd, args)

		// check catalog file path
		if !utils.FileExists(rootOptCatalogPath) && loadCatalogFatal {
			log.Fatalf("catalog not found %s", rootOptCatalogPath)
		}

		// spinner
		s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
		s.Suffix = " loading catalog..."
		err = s.Color("blue")
		if err != nil {
			log.Error(err)
		}
		s.Start()
		defer func() {
			s.Stop()
		}()

		// load tree
		loadedTree, err = tree.LoadTree(rootOptCatalogPath)
		if err != nil && loadCatalogFatal {
			log.Fatal(err)
		}
	}
}

func formatOk(selected string, treeOk bool, scriptOk bool) bool {
	var ok bool
	for _, fmt := range stringer.GetSupportedFormats(treeOk, scriptOk) {
		if fmt == selected {
			ok = true
		}
	}
	return ok
}

// Execute entry point
func Execute() error {
	return rootCmd.Execute()
}