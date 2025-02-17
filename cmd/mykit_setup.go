package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	emitter "mykit/internal/emitter/setup"
	"mykit/internal/metadata"
)

var _cmdSetup = &cobra.Command{
	Use:   "setup",
	Short: "setup",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			emitter.Setup()
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdSetup)
}
