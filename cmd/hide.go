package cmd

import (
	"github.com/spf13/cobra"
)

var hideCmd = &cobra.Command{
	Use:   "hide",
	Short: "Hide the given mods, disallowing them to be seen using 'list'",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		return modHandler.HideMods(args...)
	},
}

func init() {
	rootCmd.AddCommand(hideCmd)
}
