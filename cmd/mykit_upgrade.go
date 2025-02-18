package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/ChisTrun/trunkit/internal/constant"
	"github.com/ChisTrun/trunkit/internal/metadata"
	osutil "github.com/ChisTrun/trunkit/internal/util/os"
)

var _cmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			upgradeLatest()
			return
		}

		fmt.Println("- Current version:", metadata.MyKitVersion)
		upgrade(args[0])
	},
}

func init() {
	_rootCmd.AddCommand(_cmdUpgrade)
}

func upgradeLatest() bool {
	latestVersion := metadata.GetLatestVersion()
	if latestVersion == metadata.MyKitVersion {
		return true
	}
	fmt.Println("- Current version:", metadata.MyKitVersion)
	fmt.Println("- Upgrade version:", latestVersion)

	prompt := promptui.Prompt{
		Label:       "Do you want to upgrade",
		HideEntered: true,
		IsConfirm:   true,
	}
	_, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrAbort) {
			return true
		}

		fmt.Println("failed to prompt upgrade version", err)
		return false
	}

	return upgrade(latestVersion)
}

func upgrade(version string) bool {
	s := spinner.New(spinner.CharSets[33], 500*time.Millisecond)
	s.Prefix = fmt.Sprintf("Upgrade version %s ", version)
	s.FinalMSG = fmt.Sprintf("Upgrade version %s successfully!\n", version)
	s.Start()

	_, err := osutil.Exec([]string{
		fmt.Sprintf("go install %s@%s", constant.MyKitBase, version),
	})
	if err != nil {
		s.Stop()
		fmt.Println("upgrade failed", version, err)
		return false
	}

	s.Stop()
	fmt.Println()

	return true
}
