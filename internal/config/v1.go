package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
)

type ConfigV1 struct {
	BaseConfig `yaml:",inline"`
	Project    struct {
		ProjectBaseConfig `yaml:",inline"`
	} `yaml:"project"`
	Generate struct {
		GenerateBaseConfig `yaml:",inline"`
		Proto              struct {
			ProtoBaseConfig `yaml:",inline"`
			Imports         map[string]*ImportConfig `yaml:"imports"`
			Go              []string                 `yaml:"go"`
			Js              []string                 `yaml:"js"`
			JsConnect       []string                 `yaml:"js_connect"`
			Swift           []string                 `yaml:"swift"`
			Java            []string                 `yaml:"java"`
		}
	} `yaml:"generate"`
}

func NewConfigV1() *ConfigV1 {
	var config ConfigV1
	err := defaults.Set(&config)
	if err != nil {
		fmt.Println("set default config failed", err)
		os.Exit(1)
	}
	return &config
}

func (config *ConfigV1) convertToGenerateConfig() *GenerateConfig {
	genConfig := GenerateConfig{}
	genConfig.BaseConfig = config.BaseConfig
	genConfig.Project.ProjectBaseConfig = config.Project.ProjectBaseConfig
	genConfig.Generate.GenerateBaseConfig = config.Generate.GenerateBaseConfig
	genConfig.Generate.Proto.ProtoBaseConfig = config.Generate.Proto.ProtoBaseConfig
	genConfig.Generate.Proto.Go = config.Generate.Proto.Go
	genConfig.Generate.Proto.Js = map[string][]string{"": config.Generate.Proto.Js}
	genConfig.Generate.Proto.JsConnect = map[string][]string{"": config.Generate.Proto.JsConnect}
	genConfig.Generate.Proto.Swift = config.Generate.Proto.Swift
	genConfig.Generate.Proto.Java = config.Generate.Proto.Java
	genConfig.Generate.Proto.Imports = config.Generate.Proto.Imports

	return &genConfig
}
