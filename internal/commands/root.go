/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"gocatcli/internal/catalog"
	"gocatcli/internal/colorme"
	"gocatcli/internal/log"
	"gocatcli/internal/stringer"
	"gocatcli/internal/tree"
	"gocatcli/internal/utilities"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version     = "1.1.2"
	myName      = "gocatcli"
	defCatalog  = "gocatcli.catalog"
	rootTree    *tree.Tree
	rootCatalog *catalog.Catalog
	separator   = ","

	rootCmd = &cobra.Command{
		Use:     "gocatcli",
		Short:   "gocatcli - filesystem indexer",
		Long:    `The command line catalog tool for your offline data`,
		Version: version,
	}

	rootOptCatalogPath string
	rootOptDebugMode   bool
	rootOptNoColor     bool
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
	rootCmd.PersistentFlags().BoolVar(&rootOptNoColor, "nocolor", false, "disable colors")
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

		// colors
		if rootOptNoColor {
			colorme.UseColors = false
		}

		// check catalog file path
		if !utilities.FileExists(rootOptCatalogPath) && loadCatalogFatal {
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

		// load catalog
		rootCatalog = catalog.NewCatalog(rootOptCatalogPath)
		if rootCatalog == nil {
			log.Fatalf("cannot construct catalog")
		}
		rootTree, err = rootCatalog.LoadTree()
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
