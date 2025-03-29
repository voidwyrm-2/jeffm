package cmd

import (
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the given mods",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		return modHandler.UninstallMods(args...)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
