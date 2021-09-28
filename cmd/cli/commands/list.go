package commands

import (
	"github.com/spf13/cobra"
)

// init
func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of archers",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
