package build

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	osutil "mykit/internal/util/os"
)

var (
	cmd                *exec.Cmd
	default_dockerfile = "default.dockerfile"
)

//go:embed default.dockerfile
var f embed.FS

// Run "docker build ..." command
func Build(dockerfile string, registry string, repository string, tag string, context string, buildArgs []string, extraFlags string) {
	current_path, _ := os.Getwd()

	// [Step 1]: Create ".mykit_tmp" directory & load all util files into it
	tempDir := filepath.Join(current_path, ".mykit_tmp")
	if err := osutil.CreateDirIfNotExists(tempDir, 0755); err != nil {
		fmt.Println("create file dir failed", filepath.Base(tempDir), err)
		os.Exit(1)
	}
	load_file(default_dockerfile, tempDir)

	// [Step 2]: Check if "dockerfile" variable is empty -> use default dockerfile
	if dockerfile == "" {
		fmt.Println("- Custom dockerfile doesn't exist -> use default.dockerfile:")
		dockerfile = tempDir + "/" + "default.dockerfile"
	}
	// Check if the Dockerfile exists
	_, err1 := os.Stat(dockerfile)
	if err1 == nil {
		fmt.Println("- DockerFile exists:", dockerfile)
	} else if os.IsNotExist(err1) {
		fmt.Println("- DockerFile does not exist:", dockerfile)
	} else {
		fmt.Println("Error File existance:", err1)
	}

	// [Step 3]: Set default env values for cmd "go build ..."
	if os.Getenv("CGO_ENABLED") == "" {
		err := os.Setenv("CGO_ENABLED", "0")
		if err != nil {
			log.Printf("Error setting $CGO_ENABLED environment variable: %s\n", err)
			return
		}
	}

	if os.Getenv("GOOS") == "" {
		err := os.Setenv("GOOS", "linux")
		if err != nil {
			log.Printf("Error setting $GOOS environment variable: %s\n", err)
			return
		}
	}

	if os.Getenv("GOARCH") == "" {
		err := os.Setenv("GOARCH", "amd64")
		if err != nil {
			log.Printf("Error setting $GOARCH environment variable: %s\n", err)
			return
		}
	}

	// [Step 4]: Run "go build"
	cmd := exec.Command("go", "build", "-mod=readonly", "-o", "./cmd/main", "./cmd/main.go")
	cmd.Dir = current_path
	fmt.Println("- Go build command:\n", cmd.Args)
	// Set environment variables if needed
	// cmd.Env = append(os.Environ(), "KEY=VALUE")
	// Capture the command's output
	cmdGoBuild, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error Go build:", err)
		return
	}
	fmt.Println("- GO build output:", string(cmdGoBuild))

	// [Step 5]: Parse "--docker-args"
	var _buildArgs []string
	for _, arg := range buildArgs {
		_buildArgs = append(_buildArgs, "--build-arg", arg)
	}

	// [Step 6]: Parse "--extra-flags"
	var extraBuildFlags []string
	parts := strings.Split(extraFlags, " ")
	for _, arg := range parts {
		extraBuildFlags = append(extraBuildFlags, arg)
	}

	// [Step 7]: Set image name
	imageName := registry + "/" + repository + ":" + tag

	// [Step 8]: Setup "docker build"
	// Append --build-arg values
	buildCmd := append([]string{
		"docker", "build",
		"-f", dockerfile,
		"-t", imageName,
	}, _buildArgs...)

	// Append --extra-flags values``
	if extraFlags != "" {
		fmt.Println("extraFlags is not an empty string")
		buildCmd = append(buildCmd, extraBuildFlags...)
	}

	// Append build context
	buildCmd = append(buildCmd, context)

	// Append all args together
	cmd = exec.Command(buildCmd[0], buildCmd[1:]...)
	fmt.Println("- Docker build command:\n", cmd.Args)

	// [Step 9]: Set the standard output and error output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// [Step 10]: Execute "docker build" command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error Docker build:", err)
		os.RemoveAll(filepath.Join(current_path, "cmd", "main"))
		os.RemoveAll(filepath.Join(current_path, ".mykit_tmp"))
		return
	}

	// [Step 11]: Remove all temporary files & directories
	os.RemoveAll(filepath.Join(current_path, "cmd", "main"))
	os.RemoveAll(filepath.Join(current_path, ".mykit_tmp"))
}

///////////////////
// Helper Functions
///////////////////

// Load default files in mykit to temporary directory
func load_file(file string, tempDir string) {
	bytes, err := f.ReadFile(file)
	if err != nil {
		fmt.Printf("Error: Read %s fail \n", file)
		os.Exit(1)
	}

	dest := filepath.Join(tempDir, file)

	err = os.WriteFile(dest, bytes, 0755)
	if err != nil {
		fmt.Printf("Error: write %s into temp folder fail !!!", file)
		os.Exit(1)
	}

	fmt.Println("- Loaded files:\n", " +", filepath.Base(file), "->", dest)
}
