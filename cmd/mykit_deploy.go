package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	emitter "github.com/ChisTrun/trunkit/internal/emitter/deploy"
	"github.com/ChisTrun/trunkit/internal/metadata"
)

var (
	_valuesFile    string
	_extraHelmArgs string
)

var _cmdDeploy = &cobra.Command{
	Use:   "helm deploy --context [contextName] --namespace [namespace] --valuesFile [valuesFile] --service [serviceName] --tag [tag]",
	Short: "Deploy to K8s",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			if len(_context) == 0 {
				_context = os.Getenv("CONTEXT")
			}
			if len(_namespace) == 0 {
				_namespace = os.Getenv("NAMESPACE")
			}
			if len(_valuesFile) == 0 {
				_valuesFile = os.Getenv("VALUES_FILE")
			}
			if len(_service) == 0 {
				_service = os.Getenv("SERVICE")
			}

			if err := emitter.Deploy(_context, _namespace, _valuesFile, _service, _tag, _extraHelmArgs); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdDeploy)
	_cmdDeploy.Flags().StringVarP(&_context, "context", "c", "", "name of the kubeconfig context to use")
	_cmdDeploy.Flags().StringVarP(&_namespace, "namespace", "n", "", "namespace scope for this request")
	_cmdDeploy.Flags().StringVarP(&_valuesFile, "valuesFile", "f", "", "specify values in a YAML file")
	_cmdDeploy.Flags().StringVarP(&_service, "service", "s", "", "service name")
	_cmdDeploy.Flags().StringVarP(&_tag, "tag", "t", "dev", "tag name")
	_cmdDeploy.Flags().StringVarP(&_extraHelmArgs, "extra-helm-args", "x", "", "set extra values on the command line")
	_cmdDeploy.MarkFlagRequired("context")
	_cmdDeploy.MarkFlagRequired("service")
	_cmdDeploy.MarkFlagRequired("valuesFile")
	_cmdDeploy.MarkFlagRequired("namespace")
}
