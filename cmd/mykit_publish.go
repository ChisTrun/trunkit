package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"mykit/internal/config"
	"mykit/internal/emitter/publish"
	"mykit/internal/metadata"
)

var (
	_bucket     string
	_region     string
	_repo_token string
)

var _cmdPublish = &cobra.Command{
	Use:   "publish",
	Short: "publish proto to npm registry",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Mykit Current version:", metadata.MyKitVersion)

			if len(args) == 0 || args[0] == "" {
				fmt.Println("please input version to publish")
				os.Exit(1)
			}

			publish.Publish(config.Load().SetPackageVersion(args[0]))
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdPublish)
	_cmdPublish.AddCommand(_cmdPublishSwift)
	_cmdPublish.AddCommand(_cmdPublishJava)
	_cmdPublishSwift.Flags().StringVarP(&_bucket, "bucket", "b", os.Getenv("SWIFT_S3_BUCKET"), "S3 bucket to store swift packages")
	_cmdPublishSwift.Flags().StringVarP(&_region, "region", "r", os.Getenv("SWIFT_S3_REGION"), "S3 bucket's region")
	_cmdPublishJava.Flags().StringVarP(&_repo_token, "token", "t", os.Getenv("RELEASE_PROTO_REPO_TOKEN"), "Gitlab personal access token")
}

var _cmdPublishSwift = &cobra.Command{
	Use:   "swift",
	Short: "publish swift proto to s3",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			if len(args) == 0 || args[0] == "" {
				fmt.Println("please input version to publish")
				os.Exit(1)
			}

			err := publish.PublishSwift(_bucket, _region, config.Load().SetPackageVersion(args[0]))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

var _cmdPublishJava = &cobra.Command{
	Use:   "java",
	Short: "Publish java proto to gitlab package registry",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		if _skipUpgrade {
			fmt.Println("- Mykit Current version:", metadata.MyKitVersion)

			if len(args) == 0 || args[0] == "" {
				fmt.Println("Please input version to publish")
				os.Exit(1)
			}
			fmt.Printf("args[0]\n")
			err := publish.PublishJava(config.Load().SetPackageVersion(args[0]).SetRepoToken(_repo_token))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}
