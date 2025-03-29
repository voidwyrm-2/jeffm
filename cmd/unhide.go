package cmd

import (
	"github.com/spf13/cobra"
)

var unhideCmd = &cobra.Command{
	Use:   "unhide",
	Short: "Unhide the given mods, allowing them to be seen using 'list'",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		return modHandler.UnhideMods(args...)
	},
}

func init() {
	rootCmd.AddCommand(unhideCmd)
}
