package gocmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"

	"trunkit/internal/config"
	"trunkit/internal/emitter/common"
	"trunkit/internal/metadata"
	osutil "trunkit/internal/util/os"
)

func Init(pkg string, cfg *config.GenerateConfig) {
	if metadata.SkipMod {
		return
	}

	_, err := osutil.Exec([]string{
		fmt.Sprintf("cd %s", getGoModDir(cfg)),
		fmt.Sprintf("go mod init %s", pkg),
	})
	if err != nil {
		fmt.Println("go mod init failed", pkg, err.Error())
		return
	}
}

func Vendor(cfg *config.GenerateConfig) {
	if metadata.SkipMod {
		return
	}

	s := spinner.New(spinner.CharSets[33], 500*time.Millisecond)
	s.Prefix = "go mod vendor "
	s.Start()
	defer s.Stop()

	_, err := osutil.Exec([]string{
		fmt.Sprintf("cd %s", getGoModDir(cfg)),
		"go mod tidy",
		"go mod vendor",
	})
	if err != nil {
		fmt.Println("go mod vendor failed", err.Error())
		return
	}
}

func Bootstrap(packages []string, cfg *config.GenerateConfig) {
	if metadata.SkipMod {
		return
	}

	common.Render(
		"goservice/internal/z_bootstrap.go.tmpl",
		filepath.Join(metadata.Dir, "z_bootstrap.go"),
		map[string]interface{}{
			"Imports": packages,
		},
		common.Overwrite())

	Vendor(cfg)

	if err := os.Remove(filepath.Join(metadata.Dir, "z_bootstrap.go")); err != nil {
		fmt.Println("remove z_bootstrap.go failed", err.Error())
		return
	}
}

func Get(packages []string, cfg *config.GenerateConfig) {
	if metadata.SkipMod {
		return
	}

	s := spinner.New(spinner.CharSets[33], 500*time.Millisecond)
	s.Prefix = "go mod vendor "
	s.Start()
	defer s.Stop()

	for _, pkg := range packages {
		_, err := osutil.Exec([]string{
			fmt.Sprintf("cd %s", getGoModDir(cfg)),
			"go get " + pkg,
		})
		if err != nil {
			fmt.Println("go get failed", pkg, err.Error())
			return
		}
	}
}

func getGoModDir(cfg *config.GenerateConfig) string {
	if cfg.Project.Monorepo {
		return filepath.Dir(metadata.Dir)
	}

	return metadata.Dir
}
