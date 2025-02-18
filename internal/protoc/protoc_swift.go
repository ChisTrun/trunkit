package protoc

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"trunkit/internal/config"
	"trunkit/internal/metadata"
)

func Swift(cfg *config.GenerateConfig) string {
	tempVersion := time.Now().Format("20060102150405")
	tempDir := filepath.Join("tmp", tempVersion)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Println("make temp directory failed", err)
		os.Exit(1)
	}

	outPaths := []string{
		fmt.Sprintf("--swift_out=Visibility=Public:%s", filepath.Join(metadata.Dir, tempDir)),
		fmt.Sprintf("--grpc-swift_out=Visibility=Public:%s", filepath.Join(metadata.Dir, tempDir)),
	}
	executeProtoc("api", cfg.Generate.Proto.Swift, outPaths, cfg)
	return filepath.Join(tempDir, filepath.Base(metadata.Dir), "api")
}
