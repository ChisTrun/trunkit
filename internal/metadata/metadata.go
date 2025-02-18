package metadata

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/ChisTrun/trunkit/internal/constant"
	osutil "github.com/ChisTrun/trunkit/internal/util/os"
)

var (
	MyKitVersion = "unknown"
	Dir          = "."

	SkipMod bool // skip executing go mod
)

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		MyKitVersion = buildInfo.Main.Version
	}

	dir, err := os.Getwd()
	if err == nil {
		Dir = dir
	}
}

func GetLatestVersion() string {
	s := spinner.New(spinner.CharSets[33], 500*time.Millisecond)
	s.Prefix = "Check version "
	s.Start()
	defer s.Stop()

	cmd := fmt.Sprintf("git ls-remote --tags --sort=-version:refname https://%s.git | grep -o 'v.*' | head -1", constant.MyKitBase)
	latestVersion, err := osutil.Exec([]string{cmd})
	if err != nil {
		fmt.Println("get current version failed", err)
		return ""
	}

	return strings.ReplaceAll(strings.TrimSpace(latestVersion), "\n", "")
}
