package cmd

import (
	"fmt"
	"os"

	"trunkit/internal/config"
	"trunkit/internal/metadata"

	"github.com/spf13/cobra"

	emitter "trunkit/internal/emitter/generate/go"
	"trunkit/internal/emitter/migrate"
)

var _cmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			if len(args) == 0 || args[0] == "" {
				fmt.Println("please input source path")
				os.Exit(1)
			}

			migrate.Migrate(args[0], _name, config.Load())
			emitter.Generate(config.Load())
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdMigrate)
	_cmdMigrate.Flags().StringVar(&_name, "name", "", "name of the service")
}
