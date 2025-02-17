package deploy

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	osutil "mykit/internal/util/os"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed mychart-0.0.1.tgz
	f               embed.FS
	cmd             *exec.Cmd
	chartRepository = "mychart-0.0.1.tgz"
	appVersion      = "dev"
)

func embedFile(fileName string, destDir string) error {
	bytes, err := f.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(destDir, fileName), bytes, 0755); err != nil {
		return err
	}

	return nil
}

func getTagFromExtraArgs(extraHelmArgs string) string {
	key := "image.customTag"
	separator := "[,\\s]+" // Match one or more commas or spaces as a separator
	pairs := regexp.MustCompile(separator).Split(extraHelmArgs, -1)

	// Construct the regular expression to match the key and its value
	pattern := fmt.Sprintf("^%s=(.*)$", key)
	re := regexp.MustCompile(pattern)

	// Search for the key-value pair in the list of pairs
	for _, pair := range pairs {
		match := re.FindStringSubmatch(pair)
		if len(match) > 1 {
			value := match[1]
			return value
		}
	}
	return appVersion
}

func injectAppVersion(destDir string, tag string) error {
	chartPath := filepath.Join(destDir, "mychart", "Chart.yaml")
	yamlFile, err := os.ReadFile(chartPath)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if err := yaml.Unmarshal(yamlFile, &data); err != nil {
		return fmt.Errorf("error unmarshaling YAML: %v", err)
	}
	data["appVersion"] = tag
	newYamlContent, _ := yaml.Marshal(&data)
	os.WriteFile(chartPath, newYamlContent, 0644)
	return nil
}

func Deploy(context string, namespace string, valuesFile string, service string, tag string, extraHelmArgs string) error {
	currentPath, _ := os.Getwd()
	destDir := filepath.Join(currentPath, "tmp")

	defer os.RemoveAll(destDir)

	var customArgArray []string
	for _, ele := range strings.Split(extraHelmArgs, " ") {
		if ele != "--set" {
			customArgArray = append(customArgArray, ele)
		}
	}

	if err := osutil.CreateDirIfNotExists(destDir, 0755); err != nil {
		return fmt.Errorf("create file %s failed: %v", filepath.Base(destDir), err)

	}

	embedFile(chartRepository, destDir)

	cmd = exec.Command("tar", "-xvzf", filepath.Join(destDir, chartRepository), "-C", destDir)
	cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("can not decompress helm repo")
	}

	if tag != "" {
		injectAppVersion(destDir, tag)
	} else {
		injectAppVersion(destDir, getTagFromExtraArgs(extraHelmArgs))
	}

	if _, err := os.Stat(valuesFile); err != nil {
		return fmt.Errorf("%s does not exist", valuesFile)
	}

	cmd = exec.Command("helm", "upgrade", "--install", service, filepath.Join(destDir, "mychart"), "-f", valuesFile, "--set", fmt.Sprintf("image.customTag=%s", tag), "--set", strings.Join(customArgArray, ","), fmt.Sprintf("--kube-context=%s", context), fmt.Sprintf("--namespace=%s", namespace), "--force")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if cmd.ProcessState.ExitCode() != 0 {
		return err
	}

	return nil
}
