/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/jeffm/modapi"
)

var (
	version    string
	modHandler modapi.ModHandler
)

var rootCmd = &cobra.Command{
	Use:   "Rush",
	Short: "Rush is a terminal-based mod manager for Haste: Broken Worlds",
	Long:  ``,
}

func Execute(_version string) error {
	version = _version

	hh, err := modapi.NewHomeHandler()
	if err != nil {
		return err
	}

	err = hh.VerifyRushFolder()
	if err != nil {
		return err
	}

	modHandler, err = modapi.NewModHandler(hh)
	defer modHandler.Close()
	if err != nil {
		return err
	}

	err = rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

func init() {
}
