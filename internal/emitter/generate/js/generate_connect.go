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

func GenerateConnectJs(cfg *config.GenerateConfig) (outPaths []string) {
	for packagePostfix, protoFiles := range cfg.Generate.Proto.JsConnect {
		if len(protoFiles) == 0 {
			break
		}

		outPath := protoc.JsConnect(protoFiles, cfg)
		outPaths = append(outPaths, outPath)

		_replacers = getReplacersConnect(protoFiles, cfg)

		_filesCleanedInfo.isCleaningJsConnect = true
		err := filepath.Walk(outPath, cleanJSFile)
		if err != nil {
			fmt.Println("explore files failed", err.Error())
			os.Exit(1)
		}
		_filesCleanedInfo.isCleaningJsConnect = false

		exports, typesVersions := getExportsAndTypesVersions(_filesCleanedInfo.jsConnectFileCleaned)
		_filesCleanedInfo.jsConnectFileCleaned = []string{} // reset js connect file list

		dependencies, registries := getDependencies(true, cfg)
		devDependencies := getDevDependencies()
		peerDependencies := getPeerDependencies(true)

		for template, dest := range _connectTemplates {
			common.Render(
				template,
				filepath.Join(outPath, dest),
				map[string]interface{}{
					"Package":          cfg.Project.NpmPackage + packagePostfix + "-connect",
					"Name":             cfg.Project.Name,
					"Registry":         cfg.Project.NpmRegistry,
					"Namespace":        getNpmNamespace(cfg.Project.NpmPackage),
					"Dependencies":     dependencies,
					"DevDependencies":  devDependencies,
					"PeerDependencies": peerDependencies,
					"Registries":       registries,
					"Version":          cfg.PackageVersion,
					"MyKitVersion":     metadata.MyKitVersion,
					"Author":           "", // TODO consider,
					"Exports":          exports,
					"TypesVersions":    typesVersions,
				},
			)
		}

	}
	return outPaths
}

const _connectImportRegex = `"(.*)%s"`

func getReplacersConnect(protoFiles []string, cfg *config.GenerateConfig) []*replacer {
	var replacers []*replacer

	for _, fileName := range protoFiles {
		filePath := filepath.Join(metadata.Dir, "api", fileName)
		definition, err := protoc.ParseProto(filePath)
		if err != nil {
			fmt.Println("parse file failed", filePath, err)
			os.Exit(1)
		}

		var npmPackages = map[string]string{}
		for p, i := range cfg.Generate.Proto.Imports {
			npmPackages[p] = i.NpmPackage + "-connect"
		}

		replacers = append(replacers, createConnectReplacers(fileName, filepath.Join(filepath.Base(metadata.Dir), "api", fileName), npmPackages, cfg)...)

		proto.Walk(definition,
			proto.WithImport(func(i *proto.Import) {
				replacers = append(replacers, createConnectReplacers(fileName, i.Filename, npmPackages, cfg)...)

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
					regexp:   regexp.MustCompile(fmt.Sprintf(_connectImportRegex, pb)),
					newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pb)),
				})

				pbJs := strings.ReplaceAll(i.Filename, ".proto", "_pb.js")
				pbmJs := strings.ReplaceAll(i.Filename, ".proto", "_pb.mjs")

				replacers = append(replacers,
					&replacer{
						regexp:   regexp.MustCompile(fmt.Sprintf(_connectImportRegex, pbJs)),
						newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pbJs)),
					},
					&replacer{
						regexp:   regexp.MustCompile(fmt.Sprintf(_connectImportRegex, pbmJs)),
						newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pbmJs)),
					},
				)
			}),
		)
	}

	return replacers
}

func createConnectReplacers(fileName, importName string, npmPackages map[string]string, cfg *config.GenerateConfig) []*replacer {
	var replacers []*replacer

	importPrefix, found := npmPackages[importName]
	if !found {
		importPrefix = cfg.Project.NpmPackage + "-connect"
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
		regexp:   regexp.MustCompile(fmt.Sprintf(_connectImportRegex, pb)),
		newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pb)),
	})

	pbJs := strings.ReplaceAll(importName, ".proto", "_pb.js")
	replacers = append(replacers, &replacer{
		regexp:   regexp.MustCompile(fmt.Sprintf(_connectImportRegex, pbJs)),
		newValue: fmt.Sprintf("'%s/%s'", importPrefix, filepath.Base(pbJs)),
	})

	return replacers
}

func getExportsAndTypesVersions(fileNames []string) (string, string) {
	exports := []string{}
	typesVersions := []string{}

	for _, fileName := range fileNames {
		if !strings.HasSuffix(fileName, ".ts") || strings.HasSuffix(fileName, ".d.ts") {
			continue
		}
		name := fileName[:len(fileName)-3]
		exports = append(
			exports,
			fmt.Sprintf(
				`"./%s" : {
					"types": "./dist/types/%s.d.mts",
					"import": "./dist/esm/%s.mjs",
					"require": "./dist/cjs/%s.js"
				}`,
				name, name, name, name,
			),
			fmt.Sprintf(
				`"./%s.mjs" : {
					"types": "./dist/types/%s.d.mts",
					"import": "./dist/esm/%s.mjs",
					"require": "./dist/cjs/%s.js"
				}`,
				name, name, name, name,
			),
		)
		typesVersions = append(
			typesVersions,
			fmt.Sprintf(
				`"%s" : [ "./dist/types/%s.d.mts" ]`,
				name, name,
			),
			fmt.Sprintf(
				`"%s.mjs" : [ "./dist/types/%s.d.mts" ]`,
				name, name,
			),
		)
	}

	exportsStr := strings.Join(exports, ",\n\t\t")
	typesVersionsStr := fmt.Sprintf(
		`"*": {
			%s
		}`,
		strings.Join(typesVersions, ",\n\t\t\t"),
	)
	return exportsStr, typesVersionsStr
}
