package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the available mods",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mods, err := modHandler.GetMods()
		if err != nil {
			return err
		}

		fmt.Printf("%d mods currently installed:\n", len(mods))

		for _, m := range mods {
			if strings.HasPrefix(m, "[ENABLED]") {
				fmt.Println("\033[92m" + m + "\033[0m")
			} else {
				fmt.Println("\033[91m" + m + "\033[0m")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
