package cmd

import (
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables the given mods",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		return modHandler.DisableMods(args...)
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
}
