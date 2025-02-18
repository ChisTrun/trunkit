package cmd

import (
	"fmt"

	"trunkit/cmd/flags"
	emitter "trunkit/internal/emitter/build"
	"trunkit/internal/metadata"

	"github.com/spf13/cobra"
)

var _cmdBuild = &cobra.Command{
	Use:   "build --dockerfile [Dockerfile] --image-registry [GirRegistry] --image-repository [ImageRepositoy] --tag [ImageTag] --context [BuildContext] --build-arg [BuildArgument] --extra-docker-args [ExtraDockerArgs]",
	Short: "Build image for a specific service",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)
			emitter.Build(_dockerfile, _registry, _repository, _tag, _context, _buildArgs, _extraFlags)
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	flags.StringPersistentFlag(_cmdBuild, &_dockerfile, "dockerfile", "f", "", "Path to the Dockerfile (default: ${PWD}/.mykit_tmp/default.dockerfile)", "DOCKERFILE", true)
	flags.StringPersistentFlag(_cmdBuild, &_registry, "image-registry", "r", "", "Image registry -> e.g registry.ugaming.io", "IMAGE_REGISTRY", true)
	flags.StringPersistentFlag(_cmdBuild, &_repository, "image-repository", "i", "", "Image repository -> e.g marketplace/packages", "IMAGE_REPOSITORY", true)
	flags.StringPersistentFlag(_cmdBuild, &_tag, "tag", "t", "", "Name and optionally a tag (e.g. dev, stg, 1.0.0, 2.0.1, etc)", "TAG", true)
	flags.StringPersistentFlag(_cmdBuild, &_context, "context", "c", "", "Build Context", "BUILD_CONTEXT", true)
	flags.StringPersistentFlag(_cmdBuild, &_extraFlags, "extra-docker-args", "x", "", "Extra docker build flags (e.g mykit build ... -x '--no-cache --label ABC_XYZ')", "", false)
	_cmdBuild.Flags().StringSliceVarP(&_buildArgs, "build-arg", "a", []string{}, "Set build-time variables (e.g mykit build ... --build-arg 'SERVICE=abc' --build-arg 'GO_VERSION=1.19')")
	_rootCmd.AddCommand(_cmdBuild)
}
