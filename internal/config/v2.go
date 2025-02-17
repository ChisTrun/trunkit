package config

import (
	"fmt"
	"os"
	"path/filepath"

	"mykit/internal/metadata"

	"github.com/creasty/defaults"
)

type ConfigV2 struct {
	BaseConfig `yaml:",inline"`
	Project    struct {
		ProjectBaseConfig `yaml:",inline"`
	} `yaml:"project"`
	Generate struct {
		GenerateBaseConfig `yaml:",inline"`
		Proto              struct {
			ProtoBaseConfig `yaml:",inline"`
			Imports         map[string]*ImportConfigV2 `yaml:"imports"`
			FileGroups      map[string][]string        `yaml:"file_groups"`
			Go              map[string]string          `yaml:"go"`
			Js              map[string]string          `yaml:"js"`
			JsConnect       map[string]string          `yaml:"js_connect"`
			Swift           map[string]string          `yaml:"swift"`
			Java            map[string]string          `yaml:"java"`
		}
	} `yaml:"generate"`
}

type ImportConfigV2 struct {
	Path            string `yaml:"path"`
	ProtoPath       string `yaml:"proto_path"`
	GoPackage       string `yaml:"go_package"`
	NpmPackage      string `yaml:"npm_package"`
	JavaPackage     string `yaml:"java_package"`
	PackageRegistry string `yaml:"package_registry"`
	ProtoFiles      map[string]struct {
		Types []string `yaml:"types"`
	} `yaml:"proto_files"`
}

func NewConfigV2() *ConfigV2 {
	var config ConfigV2
	err := defaults.Set(&config)
	if err != nil {
		fmt.Println("set default config failed", err)
		os.Exit(1)
	}
	return &config
}

func (config *ConfigV2) convertToGenerateConfig() *GenerateConfig {
	filesGroups := config.Generate.Proto.FileGroups
	genConfig := GenerateConfig{}
	genConfig.BaseConfig = config.BaseConfig
	genConfig.Project.ProjectBaseConfig = config.Project.ProjectBaseConfig
	genConfig.Generate.GenerateBaseConfig = config.Generate.GenerateBaseConfig
	genConfig.Generate.Proto.ProtoBaseConfig = config.Generate.Proto.ProtoBaseConfig
	genConfig.Generate.Proto.Go = convertProtoV2ToSingle(filesGroups, config.Generate.Proto.Go)
	genConfig.Generate.Proto.Js = convertProtoV2ToMultiple(filesGroups, config.Generate.Proto.Js)
	genConfig.Generate.Proto.JsConnect = convertProtoV2ToMultiple(filesGroups, config.Generate.Proto.JsConnect)
	genConfig.Generate.Proto.Swift = convertProtoV2ToSingle(filesGroups, config.Generate.Proto.Swift)
	genConfig.Generate.Proto.Java = convertProtoV2ToSingle(filesGroups, config.Generate.Proto.Java)
	genConfig.Generate.Proto.Imports = config.convertImport()

	return &genConfig
}

func (config *ConfigV2) convertImport() map[string]*ImportConfig {
	configConverted := make(map[string]*ImportConfig)
	for _, configGroup := range config.Generate.Proto.Imports {
		for fileName, fileConfig := range configGroup.ProtoFiles {
			configConvertedKey := fmt.Sprintf("%s/%s", configGroup.ProtoPath, fileName)
			configConverted[configConvertedKey] = &ImportConfig{
				Path:            configGroup.Path,
				GoPackage:       configGroup.GoPackage,
				NpmPackage:      configGroup.NpmPackage,
				JavaPackage:     configGroup.JavaPackage,
				PackageRegistry: configGroup.PackageRegistry,
				NpmRegistry:     configGroup.PackageRegistry,
				MavenRegistry:   configGroup.PackageRegistry,
				Types:           fileConfig.Types,
			}
		}
	}
	return configConverted
}

func convertProtoV2ToSingle(fileGroups map[string][]string, genGroups map[string]string) []string {
	var filesToGen []string
	for groupKey := range genGroups {
		// only gen files for first group found
		filesToGen = fileGroups[groupKey]
		break
	}
	if len(filesToGen) == 0 {
		filesToGen = getAllProtoFiles()
	}
	return filesToGen
}

func convertProtoV2ToMultiple(fileGroups map[string][]string, genGroups map[string]string) map[string][]string {
	filesToGenWithGroup := make(map[string][]string)
	for groupKey, postfix := range genGroups {
		filesToGen := fileGroups[groupKey]
		if len(filesToGen) == 0 {
			filesToGen = getAllProtoFiles()
		}
		packagePostfix := ""
		if len(postfix) != 0 {
			packagePostfix = "-" + postfix
		}
		filesToGenWithGroup[packagePostfix] = filesToGen
	}
	return filesToGenWithGroup
}

func getAllProtoFiles() (protoFiles []string) {
	filePath := filepath.Join(metadata.Dir, "api")
	filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".proto" {
			protoFiles = append(protoFiles, info.Name())
		}
		return nil
	})
	return protoFiles
}
