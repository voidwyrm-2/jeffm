package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var hiddenCmd = &cobra.Command{
	Use:   "hidden",
	Short: "Shows the mods that have been hidden",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mods, err := modHandler.GetHidden()
		if err != nil {
			return err
		}

		fmt.Printf("%d mods currently hidden:\n", len(mods))

		for _, mod := range mods {
			fmt.Println(mod)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hiddenCmd)
}
