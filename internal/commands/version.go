package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sunny-b/cryptkeeper/internal/version"
)

var Version = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Prints the version number of cryptkeeper",
	Long:    "Prints the version number of cryptkeeper",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
