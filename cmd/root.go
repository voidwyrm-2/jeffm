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
	Use:   "jeffm",
	Short: "JeffM is a terminal-based mod manager for Marvel Rivals",
	Long:  ``,
}

func Execute(_version string) error {
	version = _version

	_modHandler, err := modapi.NewModHandler()
	if err != nil {
		return err
	}

	modHandler = _modHandler

	err = rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

func init() {
}
