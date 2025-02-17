package js

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"mykit/internal/config"
	"mykit/internal/emitter/common"
	"mykit/internal/metadata"
	"mykit/internal/protoc"

	"github.com/emicklei/proto"
)

const _importRegex = "'(.*)%s'"

func GenerateJs(cfg *config.GenerateConfig) (outPaths []string) {
	for packagePostfix, protoFiles := range cfg.Generate.Proto.Js {
		if len(protoFiles) == 0 {
			break
		}

		outPath := protoc.Js(protoFiles, cfg)
		outPaths = append(outPaths, outPath)

		_replacers = getReplacers(protoFiles, cfg)

		err := filepath.Walk(outPath, cleanJSFile)
		if err != nil {
			fmt.Println("explore files failed", err.Error())
			os.Exit(1)
		}

		dependencies, registries := getDependencies(false, cfg)
		devDependencies := getDevDependencies()
		peerDependencies := getPeerDependencies(false)

		for template, dest := range _templates {
			common.Render(
				template,
				filepath.Join(outPath, dest),
				map[string]interface{}{
					"Package":          cfg.Project.NpmPackage + packagePostfix,
					"Name":             cfg.Project.Name,
					"Registry":         cfg.Project.NpmRegistry,
					"Namespace":        getNpmNamespace(cfg.Project.NpmPackage),
					"Dependencies":     dependencies,
					"DevDependencies":  devDependencies,
					"PeerDependencies": peerDependencies,
					"Registries":       registries,
					"Version":          cfg.PackageVersion,
					"MyKitVersion":     metadata.MyKitVersion,
					"Author":           "", // TODO consider
				},
			)
		}
	}
	return outPaths
}

func getReplacers(protoFiles []string, cfg *config.GenerateConfig) []*replacer {
	var replacers = make([]*replacer, len(_defaultReplacers))
	copy(replacers, _defaultReplacers)

	for _, fileName := range protoFiles {
		filePath := filepath.Join(metadata.Dir, "api", fileName)
		definition, err := protoc.ParseProto(filePath)
		if err != nil {
			fmt.Println("parse file failed", filePath, err)
			os.Exit(1)
		}

		var npmPackages = map[string]string{}
		for p, i := range cfg.Generate.Proto.Imports {
			npmPackages[p] = i.NpmPackage
		}

		replacers = append(replacers, createReplacers(fileName, filepath.Join(filepath.Base(metadata.Dir), "api", fileName), npmPackages, cfg)...)

		proto.Walk(definition,
			proto.WithImport(func(i *proto.Import) {
				replacers = append(replacers, createReplacers(fileName, i.Filename, npmPackages, cfg)...)

				if i.Filename == "validate/validate.proto" ||
					strings.HasPrefix(i.Filename, "google/protobuf/") {
					return
				}

				importPrefix, found := npmPackages[i.Filename]
				if !found {
					importPrefix = cfg.Project.NpmPackage
				}
				if !strings.Contains(i.Filename, "/api/") {
					importPrefix = filepath.Join(importPrefix, filepath.Dir(i.Filename))

					fmt.Printf("WARN: %s has import %s that should has more information, ex `servicex/api/%s`. \n",
						fileName, i.Filename, filepath.Base(i.Filename))
				} else {
					importPrefix = filepath.Join(importPrefix, filepath.Dir(strings.Split(i.Filename, "/api/")[1]))
				}

				pb := strings.ReplaceAll(i.Filename, ".proto", "_pb")
				replacers = append(replacers, &replacer{
					regexp:   regexp.MustCompile(fmt.Sprintf(_importRegex, pb)),
					newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pb)),
				})

				pbJs := strings.ReplaceAll(i.Filename, ".proto", "_pb.js")
				replacers = append(replacers, &replacer{
					regexp:   regexp.MustCompile(fmt.Sprintf(_importRegex, pbJs)),
					newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pbJs)),
				})
			}),
		)
	}

	return replacers
}

func createReplacers(fileName, importName string, npmPackages map[string]string, cfg *config.GenerateConfig) []*replacer {
	var replacers []*replacer

	if importName == "validate/validate.proto" ||
		strings.HasPrefix(importName, "google/protobuf/") {
		return replacers
	}

	importPrefix, found := npmPackages[importName]
	if !found {
		importPrefix = cfg.Project.NpmPackage
	}
	if !strings.Contains(importName, "/api/") {
		importPrefix = filepath.Join(importPrefix, filepath.Dir(importName))

		fmt.Printf("WARN: %s has import %s that should has more information, ex `servicex/api/%s`. \n",
			fileName, importName, filepath.Base(importName))
	} else {
		importPrefix = filepath.Join(importPrefix, filepath.Dir(strings.Split(importName, "/api/")[1]))
	}

	pb := strings.ReplaceAll(importName, ".proto", "_pb")
	replacers = append(replacers, &replacer{
		regexp:   regexp.MustCompile(fmt.Sprintf(_importRegex, pb)),
		newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pb)),
	})

	pbJs := strings.ReplaceAll(importName, ".proto", "_pb.js")
	replacers = append(replacers, &replacer{
		regexp:   regexp.MustCompile(fmt.Sprintf(_importRegex, pbJs)),
		newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pbJs)),
	})

	return replacers
}
