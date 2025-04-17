package cmd

import (
	"github.com/spf13/cobra"
)

var disable_all *bool

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables the given mods",
	Long:  ``,
	Args:  cobra.RangeArgs(0, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		if *disable_all {
			mods, err := modHandler.GetRawMods()
			if err != nil {
				return err
			}

			return modHandler.DisableMods(mods...)
		}

		return modHandler.DisableMods(args...)
	},
}

func init() {
	disable_all = disableCmd.Flags().Bool("all", false, "Disable all mods")

	rootCmd.AddCommand(disableCmd)
}
