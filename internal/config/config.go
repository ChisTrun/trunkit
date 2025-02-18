package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ChisTrun/trunkit/internal/metadata"

	"gopkg.in/yaml.v3"
)

const _configFileName = "mykit.yaml"

var _config *GenerateConfig

type GenerateConfig struct {
	BaseConfig     `yaml:",inline"`
	PackageVersion string `yaml:",omitempty"` // for js, jsconnect, java, swift
	RepoToken      string `yaml:",omitempty"` // for java generate
	Project        struct {
		ProjectBaseConfig `yaml:",inline"`
	}
	Generate struct {
		GenerateBaseConfig `yaml:",inline"`
		Proto              struct {
			ProtoBaseConfig `yaml:",inline"`
			Imports         map[string]*ImportConfig
			Go              []string
			Js              map[string][]string // map package postfix to proto files
			JsConnect       map[string][]string // map package postfix to proto files
			Swift           []string
			Java            []string
		}
	}
}

func (cfg *GenerateConfig) SetPackageVersion(version string) *GenerateConfig {
	cfg.PackageVersion = version
	return cfg
}

func (cfg *GenerateConfig) SetRepoToken(repoToken string) *GenerateConfig {
	cfg.RepoToken = repoToken
	return cfg
}

func Load() *GenerateConfig {
	if _config == nil {
		_config = resolveDefaultValues(load(filepath.Join(metadata.Dir, _configFileName)))
	}
	return _config
}

func resolveDefaultValues(cfg *GenerateConfig) *GenerateConfig {
	if cfg.Generate.Proto.Imports == nil {
		for _, im := range cfg.Generate.Proto.Imports {
			if im.PackageRegistry == "" {
				continue
			}
			if im.NpmRegistry == "" {
				im.NpmRegistry = path.Join(im.NpmRegistry, "npm")
			}
			if im.MavenRegistry == "" {
				im.MavenRegistry = path.Join(im.MavenRegistry, "maven")
			}
		}
	}
	return cfg
}

func load(filePath string) *GenerateConfig {
	fmt.Println("Load config file:")
	fmt.Println("-", filePath)

	yamlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("read config file failed", err)
		os.Exit(1)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		fmt.Println("unmarshal config failed", err)
		os.Exit(1)
	}

	version, ok := data["version"].(int)
	if !ok {
		// default version is 1
		version = 1
	}

	var cfg *GenerateConfig

	switch version {
	case 2:
		configV2 := NewConfigV2()
		err := yaml.Unmarshal(yamlData, &configV2)
		if err != nil {
			fmt.Println("unmarshal config failed, please check your config file (version 2),", err)
			os.Exit(1)
		}
		cfg = configV2.convertToGenerateConfig()
	default:
		configV1 := NewConfigV1()
		err := yaml.Unmarshal(yamlData, &configV1)
		if err != nil {
			fmt.Println("unmarshal config failed, please check your config file (version 1),", err)
			os.Exit(1)
		}
		cfg = configV1.convertToGenerateConfig()

		if cfg.Extend != "" {
			extendCfg := load(resolveExtendPath(filePath, cfg.Extend))
			merged, err := merge(cfg, extendCfg)
			if err != nil {
				fmt.Println("merge config failed", err)
				os.Exit(1)
			}
			return merged
		}
	}

	return cfg
}

func resolveExtendPath(configPath string, extendPath string) string {
	if filepath.IsAbs(extendPath) {
		return extendPath
	}
	return filepath.Join(filepath.Dir(configPath), extendPath)
}

func merge[T any](parent *T, child *T) (*T, error) {
	pYaml, _ := yaml.Marshal(parent)
	cYaml, _ := yaml.Marshal(child)
	var master map[string]interface{}
	err := yaml.Unmarshal(pYaml, &master)
	if err != nil {
		return nil, err
	}

	var override map[string]interface{}
	err = yaml.Unmarshal(cYaml, &override)
	if err != nil {
		return nil, err
	}

	for k, v := range override {
		master[k] = v
	}

	bs, err := yaml.Marshal(master)
	if err != nil {
		return nil, err
	}

	var merged T
	err = yaml.Unmarshal(bs, &merged)
	if err != nil {
		return nil, err
	}
	return &merged, nil
}
