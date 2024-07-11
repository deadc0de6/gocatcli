/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package commands

import (
	"fmt"
	"gocatcli/internal/log"
	"gocatcli/internal/stringer"
	"gocatcli/internal/utils"

	"github.com/spf13/cobra"
)

var (
	storageCmd = &cobra.Command{
		Use:   "storage",
		Short: "Manage storage",
	}

	storageListCmd = &cobra.Command{
		Use:    "list",
		Short:  "List storages",
		PreRun: preRun(true),
		RunE:   storageList,
	}

	storageRemoveCmd = &cobra.Command{
		Use:    "rm <storage>",
		Short:  "Remove a storage and its entries from the catalog",
		Args:   cobra.ExactArgs(1),
		PreRun: preRun(true),
		RunE:   storageRemove,
	}

	storageMetaCmd = &cobra.Command{
		Use:    "meta <storage-name> <meta>",
		Short:  "Update storage meta data",
		Args:   cobra.ExactArgs(2),
		PreRun: preRun(true),
		RunE:   storageMeta,
	}

	storageTagCmd = &cobra.Command{
		Use:    "tag <storage-name> <tag>",
		Short:  "Add tag to storage",
		Args:   cobra.ExactArgs(2),
		PreRun: preRun(true),
		RunE:   storageTag,
	}

	storageUntagCmd = &cobra.Command{
		Use:    "untag <storage-name> <tag>",
		Short:  "Remove tag from storage",
		Args:   cobra.ExactArgs(2),
		PreRun: preRun(true),
		RunE:   storageUntag,
	}

	storageOptIndent  bool
	storageRmOptForce bool
)

func init() {
	storageCmd.AddCommand(storageRemoveCmd)
	storageCmd.AddCommand(storageMetaCmd)
	storageCmd.AddCommand(storageTagCmd)
	storageCmd.AddCommand(storageUntagCmd)
	storageCmd.AddCommand(storageListCmd)

	rootCmd.AddCommand(storageCmd)

	// everything that affects tree gets the indent option
	storageRemoveCmd.PersistentFlags().BoolVarP(&storageOptIndent, "indent", "I", true, "do not indent json")
	storageMetaCmd.PersistentFlags().BoolVarP(&storageOptIndent, "indent", "I", true, "do not indent json")
	storageTagCmd.PersistentFlags().BoolVarP(&storageOptIndent, "indent", "I", true, "do not indent json")
	storageUntagCmd.PersistentFlags().BoolVarP(&storageOptIndent, "indent", "I", true, "do not indent json")

	// rm options
	storageRemoveCmd.PersistentFlags().BoolVarP(&storageRmOptForce, "force", "f", false, "do not ask user")
}

func storageSave() error {
	return loadedTree.Save(rootOptCatalogPath, storageOptIndent)
}

func storageRemove(_ *cobra.Command, args []string) error {
	name := args[0]
	if !storageRmOptForce && !utils.AskUser(fmt.Sprintf("Do you really want to remove storage \"%s\" and its children?", name)) {
		log.Fatal(fmt.Errorf("user interrupted"))
	}

	loadedTree.RemoveStorage(name)
	ret := storageSave()
	listStorages()
	return ret
}

func storageMeta(_ *cobra.Command, args []string) error {
	name := args[0]
	meta := args[1]

	storage := loadedTree.GetStorageByName(name)
	if storage == nil {
		return fmt.Errorf("no such storage %s", name)
	}

	storage.SetMeta(meta)
	ret := storageSave()
	listStorages()
	return ret
}

func storageTag(_ *cobra.Command, args []string) error {
	name := args[0]
	tag := args[1]

	storage := loadedTree.GetStorageByName(name)
	if storage == nil {
		return fmt.Errorf("no such storage %s", name)
	}

	storage.Tag(tag)
	ret := storageSave()
	listStorages()
	return ret
}

func storageUntag(_ *cobra.Command, args []string) error {
	name := args[0]
	tag := args[1]

	storage := loadedTree.GetStorageByName(name)
	if storage == nil {
		return fmt.Errorf("no such storage %s", name)
	}

	storage.Untag(tag)
	ret := storageSave()
	listStorages()
	return ret
}

func listStorages() {
	storages := loadedTree.GetStorages()
	if storages == nil {
		return
	}

	// get a stringer to print found nodes
	m := &stringer.PrintMode{
		FullPath:    true,
		Long:        true,
		InlineColor: false,
		RawSize:     false,
		Separator:   separator,
	}
	stringGetter, err := stringer.GetStringer(loadedTree, stringer.FormatNative, m)
	if err != nil {
		return
	}

	for _, sto := range storages {
		stringGetter.Print(sto, 0)
	}
}

func storageList(_ *cobra.Command, _ []string) error {
	listStorages()
	return nil
}
