package profile

import (
	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/jeffm/modapi"
)

var modHandler *modapi.ModHandler

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Marvel Rivals mod profiles",
}

func Profile(handler *modapi.ModHandler) *cobra.Command {
	modHandler = handler
	return profileCmd
}
