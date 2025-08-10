package js

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ChisTrun/trunkit/internal/config"
)

var _templates = map[string]string{
	"npm/package.json.tmpl":  "package.json",
	"npm/npmrc.tmpl":         ".npmrc",
	"npm/tsconfig.json.tmpl": "tsconfig.json",
	"npm/README.tmpl":        "README",
}
var _connectTemplates = map[string]string{
	"npm-connect/package.json.tmpl":      "package.json",
	"npm-connect/npmrc.tmpl":             ".npmrc",
	"npm-connect/tsconfig-cjs.json.tmpl": "tsconfig-cjs.json",
	"npm-connect/README.tmpl":            "README",
}

var _replacers = _defaultReplacers
var _filesCleanedInfo = struct {
	isCleaningJsConnect  bool
	jsConnectFileCleaned []string
}{}

type npmRegistry struct {
	Namespace string
	Registry  string
	AuthToken string
}

func getDevDependencies() string {
	devDependencies := []string{"\"typescript\": \"^4.1.3\""}
	return strings.Join(devDependencies, ",\n")
}

func getPeerDependencies(isConnectJs bool) string {
	peerDependencies := []string{}
	if isConnectJs {
		peerDependencies = []string{"\"@bufbuild/protobuf\": \"^1.2.1\""}
	} else {
		peerDependencies = []string{"\"google-protobuf\": \"^3.0.1\""}
	}
	return strings.Join(peerDependencies, ",\n")
}

func getDependencies(isConnectJs bool, cfg *config.GenerateConfig) (string, []*npmRegistry) {
	// default values
	dependencies := []string{}

	var (
		thisNamespace = getNpmNamespace(cfg.Project.NpmPackage)
		thisRegistry  = cfg.Project.NpmRegistry

		dependencyMap = map[string]struct{}{}
		registryMap   = map[string]string{
			thisNamespace: thisRegistry,
		}
	)
	for _, i := range cfg.Generate.Proto.Imports {
		if len(i.NpmPackage) == 0 {
			continue
		}

		if _, found := dependencyMap[i.NpmPackage]; !found {
			dependencyMap[i.NpmPackage] = struct{}{}
			if isConnectJs {
				dependencies = append(dependencies, fmt.Sprintf("\"%s\": \"^1.0.0\"", i.NpmPackage+"-connect"))
			} else {
				dependencies = append(dependencies, fmt.Sprintf("\"%s\": \"^1.0.0\"", i.NpmPackage))
			}

		}

		namespace := getNpmNamespace(i.NpmPackage)
		if len(registryMap[namespace]) == 0 {
			registryMap[namespace] = i.NpmRegistry
		}
	}

	var registries []*npmRegistry
	for namespace, registry := range registryMap {
		if len(registry) == 0 {
			registry = thisRegistry
		}
		registries = append(registries, &npmRegistry{
			Namespace: namespace,
			Registry:  registry,
			AuthToken: getNpmAuthToken(namespace),
		})
	}

	return strings.Join(dependencies, ",\n"), registries
}

func getNpmNamespace(npmPackage string) string {
	return strings.Split(strings.TrimPrefix(npmPackage, "@"), "/")[0]
}

func getNpmAuthToken(namespace string) string {
	// return fmt.Sprintf("${%s_NPM_TOKEN}", namespace)
	return "${BOT_PRIVATE_TOKEN}"
}

func cleanJSFile(filePath string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if fileInfo.IsDir() {
		return nil
	}

	if strings.HasSuffix(fileInfo.Name(), ".js") ||
		strings.HasSuffix(fileInfo.Name(), ".ts") {
		if _filesCleanedInfo.isCleaningJsConnect {
			_filesCleanedInfo.jsConnectFileCleaned = append(_filesCleanedInfo.jsConnectFileCleaned, fileInfo.Name())
		}

		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("read file failed", filePath, err)
			os.Exit(1)
		}

		newContent := string(fileContent)
		for _, replacer := range _replacers {
			newContent = replacer.regexp.ReplaceAllString(newContent, replacer.newValue)
		}

		err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			fmt.Println("write file failed", filePath, err)
			os.Exit(1)
		}
	}

	return nil
}

type replacer struct {
	regexp   *regexp.Regexp
	newValue string
}

var _defaultReplacers = []*replacer{
	{
		regexp: regexp.MustCompile(`import \* as validate_validate_pb from '.*/validate/validate_pb';`),
	},
	{
		regexp: regexp.MustCompile(`var validate_validate_pb = require\('.*/validate/validate_pb.js'\);`),
	},
	{
		regexp: regexp.MustCompile(`goog.object.extend\(proto, validate_validate_pb\);`),
	},
}
