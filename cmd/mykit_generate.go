package cmd

import (
	"fmt"
	"os"

	"mykit/internal/config"
	"mykit/internal/metadata"

	"github.com/spf13/cobra"

	goemitter "mykit/internal/emitter/generate/go"
	"mykit/internal/emitter/generate/java"
	jsemitter "mykit/internal/emitter/generate/js"
	"mykit/internal/emitter/generate/swift"
)

var _cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "gen",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			goemitter.Generate(config.Load())
			jsemitter.GenerateJs(config.Load().SetPackageVersion("devel"))
			jsemitter.GenerateConnectJs(config.Load().SetPackageVersion("devel"))

			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdGenerate)
	_cmdGenerate.AddCommand(_cmdGenerateGo)
	_cmdGenerate.AddCommand(_cmdGenerateJs)
	_cmdGenerate.AddCommand(_cmdGenerateConnectJs)
	_cmdGenerate.AddCommand(_cmdGenerateSwift)
	_cmdGenerate.AddCommand(_cmdGenerateJava)
	_cmdGenerate.PersistentFlags().BoolVarP(&metadata.SkipMod, "mod", "m", false, "skip executing go mod vendor")
}

var _cmdGenerateGo = &cobra.Command{
	Use:   "go",
	Short: "go",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			goemitter.Generate(config.Load())
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

var _cmdGenerateJs = &cobra.Command{
	Use:   "js",
	Short: "js",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			jsemitter.GenerateJs(config.Load().SetPackageVersion("devel"))
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

var _cmdGenerateConnectJs = &cobra.Command{
	Use:   "connect",
	Short: "connect",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			jsemitter.GenerateConnectJs(config.Load().SetPackageVersion("devel"))
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

var _cmdGenerateSwift = &cobra.Command{
	Use:   "swift",
	Short: "swift",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			generatedDir := swift.Generate(config.Load().SetPackageVersion("devel"))
			if getDefaultDir() != "" {
				fmt.Printf("Generated swift to %s\n", generatedDir)
			}
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

var _cmdGenerateJava = &cobra.Command{
	Use:   "java",
	Short: "Generate Java Protobuf and Grpc code",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			generatedDir, err := java.Generate(config.Load().SetPackageVersion("devel").SetRepoToken(""))
			if err != nil {
				fmt.Printf("Generate Java protobuf fail !!! \n")
				os.Exit(1)
			}

			if getDefaultDir() != "" {
				fmt.Printf("Generated .java to %s\n", generatedDir)
			}
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}
