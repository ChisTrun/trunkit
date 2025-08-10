package protoc

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ChisTrun/trunkit/internal/config"
	"github.com/ChisTrun/trunkit/internal/metadata"
)

func Java(cfg *config.GenerateConfig) string {
	tempVersion := time.Now().Format("20060102150405")
	tempDir := filepath.Join("tmp", tempVersion)
	if err := os.MkdirAll(filepath.Join(tempDir, "src/main/java"), 0755); err != nil {
		fmt.Println("make temp directory failed", err)
		os.Exit(1)
	}
	outPaths := []string{
		fmt.Sprintf("--java_out=%s", filepath.Join(metadata.Dir, tempDir, "src/main/java")),
		fmt.Sprintf("--grpc-java_out=%s", filepath.Join(metadata.Dir, tempDir, "src/main/java")),
	}
	executeProtoc("api", cfg.Generate.Proto.Java, outPaths, cfg)
	return tempDir
}
