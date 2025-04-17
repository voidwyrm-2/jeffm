package cmd

import (
	"github.com/spf13/cobra"
)

var uninstall_all *bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the given mods",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		if *uninstall_all {
			mods, err := modHandler.GetRawMods()
			if err != nil {
				return err
			}

			return modHandler.UninstallMods(mods...)
		}

		return modHandler.UninstallMods(args...)
	},
}

func init() {
	uninstall_all = uninstallCmd.Flags().Bool("all", false, "Uninstall all mods")

	rootCmd.AddCommand(uninstallCmd)
}
