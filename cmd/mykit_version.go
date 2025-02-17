package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"mykit/internal/metadata"
)

var _cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "get version",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(metadata.MyKitVersion)
	},
}

func init() {
	_rootCmd.AddCommand(_cmdVersion)
}
