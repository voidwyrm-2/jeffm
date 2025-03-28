package cmd

import (
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables the given mods",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		return modHandler.EnableMods(args...)
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
}
