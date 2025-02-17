package java

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mykit/internal/config"
	"mykit/internal/emitter/common"
	"mykit/internal/metadata"
	"mykit/internal/protoc"
)

// Generate Java protobuf -> return(
//
// releaseProtoDir: generated protobuf directory & java proto template,
//
// error: error status
// )
func Generate(cfg *config.GenerateConfig) (string, error) {

	// Step 1: Generate java protobuf
	fmt.Printf("Generate proto for java - version %s\n", cfg.PackageVersion)
	releaseProtoDir := protoc.Java(cfg)
	fmt.Printf("Generated java proto files in -> %s\n", releaseProtoDir)

	dependencies, registries := getDependencies(cfg)

	// Step 2: Create Java proto project & copy all generated proto files to it if need
	if cfg.Generate.Proto.Java != nil && cfg.RepoToken != "" {
		fmt.Printf("Generate java template project in -> %s\n", releaseProtoDir)
		err := generateJavaTemplateHelper(releaseProtoDir, dependencies, registries, cfg)
		if err != nil {
			return releaseProtoDir, err
		}
		return releaseProtoDir, nil
	}
	return releaseProtoDir, nil
}

//////////////////
// HELPER FUNCTIONS
//////////////////

var templates = map[string]string{
	"java-proto-project/build.gradle.tmpl":    "build.gradle",
	"java-proto-project/settings.gradle.tmpl": "settings.gradle",
}

// Create Java Project Template -> return (error status)
func generateJavaTemplateHelper(releaseProtoDir string, dependencies, registries []string, cfg *config.GenerateConfig) error {
	// Check if $RELEASE_PROTO_REPO_TOKEN is unset, get value from --token flag
	gitlabToken := os.Getenv("RELEASE_PROTO_REPO_TOKEN")
	if gitlabToken == "" && len(cfg.RepoToken) > 0 {
		gitlabToken = cfg.RepoToken
	}

	for _template, dest := range templates {
		common.Render(
			_template,
			filepath.Join(releaseProtoDir, dest),
			map[string]interface{}{
				"ServiceName":     cfg.Project.Name,
				"Namespace":       strings.ReplaceAll(cfg.Project.Namespace, "/", "."),
				"Version":         cfg.PackageVersion,
				"MavenRegistry":   cfg.Project.MavenRegistry,
				"MavenRegistries": registries,
				"Dependencies":    dependencies,
				"RepoToken":       gitlabToken,
				"MyKitVersion":    metadata.MyKitVersion,
			},
		)
	}

	return nil
}

func getDependencies(cfg *config.GenerateConfig) ([]string, []string) {
	var dependencyMap = make(map[string]bool)
	var registryMap = make(map[string]bool)
	if cfg.Project.MavenRegistry != "" {
		registryMap[cfg.Project.MavenRegistry] = true
	}
	for _, i := range cfg.Generate.Proto.Imports {
		if i.JavaPackage != "" && i.MavenRegistry != "" {
			dependencyMap[i.JavaPackage] = true
			registryMap[i.MavenRegistry] = true
		}
	}
	registries, dependencies := []string{}, []string{}
	for k := range registryMap {
		registries = append(registries, k)
	}
	for k := range dependencyMap {
		dependency := fmt.Sprintf("%s:+", k)
		dependencies = append(dependencies, dependency)
	}
	return dependencies, registries
}
