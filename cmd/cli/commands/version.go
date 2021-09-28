package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// init
func init() {
	rootCmd.AddCommand(versionCmd)
}

const version = "v0.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print utility version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
