package setup

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"mykit/internal/constant"
	osutil "mykit/internal/util/os"
)

const (
	_protobufHome                = "https://github.com/protocolbuffers/protobuf/releases/tag/v3.19.4"
	_grpcWebHome                 = "https://github.com/grpc/grpc-web/releases/tag/1.3.1"
	_protocGenGotagHome          = "https://gitlab.ugaming.io/marketplace/protoc-gen-gotag"
	_cmdInstallProtocGenGo       = "go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1"
	_cmdInstallProtocGenGoGrpc   = "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2"
	_cmdInstallProtocGenValidate = "go install github.com/envoyproxy/protoc-gen-validate@v0.10.1"
	_cmdInstallProtocGrpcGateway = `go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.14.0 \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.14.0`
)

//go:embed *
var f embed.FS

func Setup() {
	fmt.Println("- OS:", runtime.GOOS)
	fmt.Println("- ARCH:", runtime.GOARCH)
	fmt.Println()

	fmt.Println("Install protoc")
	setupProtoc()

	fmt.Println("Install protoc-gen-go")
	setupProtocGenGo()

	fmt.Println("Install protoc-gen-go-grpc")
	setupProtocGenGoGrpc()

	fmt.Println("Install protoc-gen-grpc-web")
	setupProtocGenGrpcWeb()

	fmt.Println("Install protoc-gen-validate")
	setupProtocGenValidate()

	fmt.Println("Install protoc-gen-gotag")
	setupProtocGenGotag()

	fmt.Println("Install protoc-grpc-gateway")
	setupProtocGenGrpcGateway()

	fmt.Println("Install include")
	setupInclude()
}

func setupProtoc() {
	if !_supported {
		fmt.Printf("!!! your os (%s %s) is not supported, please go to %s to find your os version",
			runtime.GOOS, runtime.GOARCH, _protobufHome)
		return
	}

	writeFile(filepath.Join("prerequisites", "protoc", _goOs, _goArch, "protoc"), constant.GoBin)
}

func setupProtocGenGo() {
	_, err := osutil.Exec([]string{_cmdInstallProtocGenGo})
	if err != nil {
		fmt.Println("install protoc-gen-go failed", err)
	}
}

func setupProtocGenGoGrpc() {
	_, err := osutil.Exec([]string{_cmdInstallProtocGenGoGrpc})
	if err != nil {
		fmt.Println("install protoc-gen-go-grpc failed", err)
	}
}

func setupProtocGenGrpcWeb() {
	if !_supported {
		fmt.Printf("!!! your os (%s %s) is not supported, please go to %s to find your os version",
			runtime.GOOS, runtime.GOARCH, _grpcWebHome)
		return
	}

	writeFile(filepath.Join("prerequisites", "protoc", _goOs, _goArch, "protoc-gen-grpc-web"), constant.GoBin)
}

func setupProtocGenValidate() {
	_, err := osutil.Exec([]string{_cmdInstallProtocGenValidate})
	if err != nil {
		fmt.Println("install protoc-gen-validate failed", err)
	}
}

func setupProtocGenGotag() {
	if !_supported {
		fmt.Printf("!!! your os (%s %s) is not supported, please go to %s and rebuild for your os version",
			runtime.GOOS, runtime.GOARCH, _protocGenGotagHome)
		return
	}

	writeFile(filepath.Join("prerequisites", "protoc", _goOs, _goArch, "protoc-gen-gotag"), constant.GoBin)
}

func setupProtocGenGrpcGateway() {
	_, err := osutil.Exec([]string{_cmdInstallProtocGrpcGateway})
	if err != nil {
		fmt.Println("install protoc-gen-grpc-gateway failed", err)
	}
}

func setupInclude() {
	writeDir("prerequisites/include", constant.GoInclude)
}

func writeDir(name string, destDir string) {
	files, err := f.ReadDir(name)
	if err != nil {
		fmt.Println("read dir failed", name, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			writeDir(filepath.Join(name, file.Name()), filepath.Join(destDir, file.Name()))
			continue
		}

		writeFile(filepath.Join(name, file.Name()), destDir)
	}
}

func writeFile(filePath string, destDir string) {
	bytes, err := f.ReadFile(filePath)
	if err != nil {
		fmt.Println("read file failed", filepath.Base(filePath), err)
		return
	}

	dest := filepath.Join(destDir, filepath.Base(filePath))
	if err := osutil.CreateDirIfNotExists(destDir, 0755); err != nil {
		fmt.Println("create file dir failed", filepath.Base(filePath), err)
		return
	}
	err = os.WriteFile(dest, bytes, 0755)
	if err != nil {
		fmt.Println("write file failed", filepath.Base(filePath), err)
		return
	}

	fmt.Println(" +", filepath.Base(filePath), "->", dest)
}

var (
	_supported     bool
	_goOs, _goArch string
)

func init() {
	switch runtime.GOOS {
	case "darwin":
		_goOs = runtime.GOOS
	}

	switch runtime.GOARCH {
	case "arm64", "amd64":
		_goArch = runtime.GOARCH
	}

	switch filepath.Join(_goOs, _goArch) {
	case filepath.Join("darwin", "arm64"),
		filepath.Join("darwin", "amd64"):
		_supported = true
	}
}
