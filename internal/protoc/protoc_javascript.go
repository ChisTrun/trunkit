package protoc

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ChisTrun/trunkit/internal/config"
	"github.com/ChisTrun/trunkit/internal/metadata"
)

func Js(protoFiles []string, cfg *config.GenerateConfig) string {
	tempVersion := time.Now().Format("20060102150405")
	tempDir := filepath.Join("tmp", tempVersion)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Println("make temp directory failed", err)
		os.Exit(1)
	}

	outPaths := []string{
		fmt.Sprintf("--js_out=import_style=commonjs:%s", filepath.Join(metadata.Dir, tempDir)),
		fmt.Sprintf("--grpc-web_out=import_style=typescript,mode=grpcweb:%s", filepath.Join(metadata.Dir, tempDir)),
	}

	fmt.Println("Generate proto files for JS:")
	executeProtoc("api", protoFiles, outPaths, cfg)

	return filepath.Join(tempDir, filepath.Base(metadata.Dir), "api")
}

func JsConnect(protoFiles []string, cfg *config.GenerateConfig) string {
	tempVersion := time.Now().Format("20060102150405") + "_connect"
	tempDir := filepath.Join("tmp", tempVersion)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Println("make temp directory failed", err)
		os.Exit(1)
	}

	outPaths := []string{
		fmt.Sprintf("--es_out %s", filepath.Join(metadata.Dir, tempDir)),
		fmt.Sprintf("--connect-es_out %s", filepath.Join(metadata.Dir, tempDir)),
		"--connect-es_opt target=ts+js+dts,import_extension=.mjs",
		"--es_opt target=ts+js+dts,import_extension=.mjs",
	}

	fmt.Println("Generate proto files for JS (connect-web):")
	executeProtoc("api", protoFiles, outPaths, cfg)

	return filepath.Join(tempDir, filepath.Base(metadata.Dir), "api")
}
