package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	_rootCmd     = &cobra.Command{Use: "trunkit"}
	_skipUpgrade = false
)

func Run() {
	_rootCmd.PersistentFlags().BoolVarP(&_skipUpgrade, "skip-upgrade", "u", false, "skip upgrade")
	if err := _rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	runCmd := []string{cmd.CommandPath()}
	runCmd = append(runCmd, args...)

	var alreadySkipUpgrade bool
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if len(flag.Shorthand) > 0 {
			runCmd = append(runCmd, "-"+flag.Shorthand, flag.Value.String())
		} else {
			runCmd = append(runCmd, "--"+flag.Name, flag.Value.String())
		}

		if flag.Name == "skip-upgrade" {
			alreadySkipUpgrade = true
		}
	})
	if !alreadySkipUpgrade {
		runCmd = append(runCmd, "-u")
	}

	execCommand := exec.Command("/bin/sh", "-c", strings.Join(runCmd, " "))
	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr
	_ = execCommand.Run()
}
