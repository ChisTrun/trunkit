package publish

import (
	"time"

	"mykit/internal/config"

	"github.com/briandowns/spinner"
)

type DependencyPackage struct {
	Name      string `yaml:"package"`
	Version   string `yaml:"version"`
	Namespace string `yaml:"namespace"`
}

type PackageInfo struct {
	MykitVersion string              `yaml:"mykit_version"`
	Name         string              `yaml:"package"`
	Version      string              `yaml:"version"`
	Namespace    string              `yaml:"namespace"`
	Dependencies []DependencyPackage `yaml:"dependencies"`
}

func Publish(cfg *config.GenerateConfig) {
	npmCommands := []string{
		"npm install",
		"npm run build",
		"npm publish",
	}
	s := spinner.New(spinner.CharSets[33], 500*time.Millisecond)
	defer s.Stop()

	publishJs(cfg, npmCommands, s)
	publishConnectJs(cfg, npmCommands, s)
}
