package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"mykit/internal/constant"
	emitter "mykit/internal/emitter/init"
	"mykit/internal/metadata"
)

var (
	_dir      string
	_package  string
	_name     string
	_monorepo bool
	_type     string
)

var _cmdCreate = &cobra.Command{
	Use:   "init --dir [projectDir] --package [package] --name [projectName] --monorepo [monorepo] --type [type]",
	Short: "Initialize a new project",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _type != "service" && _type != "library" {
			fmt.Println("type must be service or library")
			os.Exit(1)
		}
		if _skipUpgrade {
			fmt.Println("- Current version:", metadata.MyKitVersion)

			if len(_name) == 0 {
				fmt.Println("Please input new project name")
				os.Exit(0)
			}

			metadata.Dir = _dir
			if len(metadata.Dir) == 0 {
				metadata.Dir = getDefaultDir()
			}
			if len(_package) == 0 {
				_package = getDefaultPackage()
			}

			emitter.Init(_package, _name, _monorepo, _type)
			return
		}

		if upgradeLatest() {
			run(cmd, args)
		}
	},
}

func init() {
	_rootCmd.AddCommand(_cmdCreate)
	_cmdCreate.Flags().StringVar(&_dir, "dir", "", "project directory")
	_cmdCreate.Flags().StringVar(&_package, "package", "", "package of the service")
	_cmdCreate.Flags().StringVar(&_name, "name", "", "name of the service")
	_cmdCreate.Flags().BoolVarP(&_monorepo, "monorepo", "", false, "project is monorepo")
	_cmdCreate.Flags().StringVar(&_type, "type", "service", "type of the service: service or library")
}

func getDefaultDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("get current dir failed", err)
		os.Exit(0)
	}

	return path.Join(currentDir, _name)
}

func getDefaultPackage() string {
	if !strings.HasPrefix(metadata.Dir, constant.GoPath) {
		return _name
	}
	pkg, err := filepath.Rel(constant.GoPath, metadata.Dir)
	if err != nil {
		fmt.Println("get default package failed", err)
		return _name
	}

	return pkg
}
