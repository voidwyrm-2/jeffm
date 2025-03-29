package cmd

import (
	"github.com/spf13/cobra"
)

var install_dontPrintInstalls *bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the given mods",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		return modHandler.InstallMods(false, !*install_dontPrintInstalls, args...)
	},
}

func init() {
	install_dontPrintInstalls = installCmd.Flags().Bool("np", false, "Don't print the files as they're installed")

	rootCmd.AddCommand(installCmd)
}
