package dockercompose

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChisTrun/trunkit/internal/config"
	"github.com/ChisTrun/trunkit/internal/emitter/common"
	osutil "github.com/ChisTrun/trunkit/internal/util/os"
)

//go:embed *
var f embed.FS

const (
	_parentDir          = "affiliate/"
	_gatewayServiceName = "gateway"
)

func Generate(addGateway bool, cfg *config.GenerateConfig) {
	fmt.Println("Generate Docker Compose:")

	fmt.Println("- Add gateway:", addGateway)
	fmt.Println()

	generateDockerCompose(addGateway, cfg)

	if addGateway {
		setupEnvoyFiles()
		generateGatewayService(cfg)
	}
}

const (
	_dockerComposeTemplate = "dockercompose/docker-compose.yaml.tmpl"
	_dockerComposeDest     = "docker-compose.yaml"
	_dockerIgnoreTemplate  = "dockercompose/dockerignore.tmpl"
	_dockerIgnoreDest      = ".dockerignore"
)

func generateDockerCompose(addGateway bool, cfg *config.GenerateConfig) {
	common.Render(
		_dockerComposeTemplate,
		_dockerComposeDest,
		map[string]interface{}{
			"ProjectName":        cfg.Project.Name,
			"ParentDir":          _parentDir,
			"AddGateway":         addGateway,
			"GatewayServiceName": _gatewayServiceName,
		},
	)

	common.Render(
		_dockerIgnoreTemplate,
		_dockerIgnoreDest,
		map[string]interface{}{},
	)
}

const (
	_gitKeepTemplate          = "gateway/descriptors/gitkeep.tmpl"
	_gitKeepDest              = _parentDir + "gateway/descriptors/.gitkeep"
	_dockerEntrypointTemplate = "gateway/scripts/docker-entrypoint.sh.tmpl"
	_dockerEntrypointDest     = _parentDir + "gateway/scripts/docker-entrypoint.sh"
	_envoyTemplate            = "gateway/scripts/envoy.yaml.tmpl"
	_envoyDest                = _parentDir + "gateway/scripts/envoy.yaml"
)

func generateGatewayService(cfg *config.GenerateConfig) {
	common.Render(
		_gitKeepTemplate,
		_gitKeepDest,
		map[string]interface{}{},
	)

	common.Render(
		_dockerEntrypointTemplate,
		_dockerEntrypointDest,
		map[string]interface{}{},
	)

	common.Render(
		_envoyTemplate,
		_envoyDest,
		map[string]interface{}{
			"ProjectName":        cfg.Project.Name,
			"GatewayServiceName": _gatewayServiceName,
		},
	)
}

func setupEnvoyFiles() {
	writeDir("envoy", _parentDir)
}

func writeDir(name string, destDir string) {
	files, err := f.ReadDir(name)
	if err != nil {
		fmt.Println("read dir failed", name, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			writeDir(filepath.Join(name, file.Name()), filepath.Join(destDir, file.Name()))
			continue
		}

		writeFile(filepath.Join(name, file.Name()), destDir)
	}
}

func writeFile(filePath string, destDir string) {
	bytes, err := f.ReadFile(filePath)
	if err != nil {
		fmt.Println("read file failed", filepath.Base(filePath), err)
		return
	}

	dest := filepath.Join(destDir, filepath.Base(filePath))
	if err := osutil.CreateDirIfNotExists(destDir, 0755); err != nil {
		fmt.Println("create file dir failed", filepath.Base(filePath), err)
		return
	}
	err = os.WriteFile(dest, bytes, 0755)
	if err != nil {
		fmt.Println("write file failed", filepath.Base(filePath), err)
		return
	}

	fmt.Println(" +", filepath.Base(filePath), "->", dest)
}
