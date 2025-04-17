package profile

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
		profiles, err := modHandler.GetProfiles()
		if err != nil {
			return err
		}

		fmt.Printf("%d profiles available:\n", len(profiles))

		fmt.Println(strings.Join(profiles, "\n"))

		return nil
	},
}

func init() {
	profileCmd.AddCommand(listCmd)
}
