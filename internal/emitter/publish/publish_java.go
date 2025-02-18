package publish

import (
	"fmt"
	"os"
	"os/exec"

	"trunkit/internal/config"
	"trunkit/internal/emitter/generate/java"
)

// Publish maven package in Gitlab
func PublishJava(cfg *config.GenerateConfig) error {
	if cfg.RepoToken == "" {
		return fmt.Errorf("[Err] repo_token is empty - Please set $RELEASE_PROTO_REPO_TOKEN variable or set Personal Access Token value through --token flag !!!")
	}

	// Step 1: Generate protobuf classes into tmp directory
	javaProjectDir, err := java.Generate(cfg)
	if err != nil {
		fmt.Printf("[Err] Generate java proto failed: %s !!!\n", err)
	}
	// Step 2: Build then Publish package into gitlab by using gradle
	err = UploadMavenPackage(javaProjectDir)
	if err != nil {
		fmt.Println("[Err] Error occur when uploading maven packge")
		return fmt.Errorf("[Err] Publish maven package failed !!!")
	}
	return nil
}

func UploadMavenPackage(javaProjectDir string) error {
	err := runCommand(javaProjectDir, "gradle", "publish")
	if err != nil {
		return err
	}

	return nil
}

func runCommand(javaProjectDir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = javaProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Execute %s %s -> in %s\n", command, args, javaProjectDir)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
