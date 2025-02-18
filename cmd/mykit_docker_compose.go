package cmd

import (
	"fmt"

	"trunkit/internal/config"
	"trunkit/internal/metadata"

	"github.com/spf13/cobra"

	dcemitter "trunkit/internal/emitter/dockercompose"
)

var _addGateway bool

var _cmdGenerateDockerCompose = &cobra.Command{
	Use:   "docker-compose",
	Short: "generate docker-compose",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			dcemitter.Generate(_addGateway, config.Load())
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdGenerateDockerCompose)

	_cmdGenerateDockerCompose.Flags().BoolVar(&_addGateway, "add-gateway", false, "add gateway service")
}
