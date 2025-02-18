package swift

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ChisTrun/trunkit/internal/config"
	"github.com/ChisTrun/trunkit/internal/protoc"
)

func Generate(cfg *config.GenerateConfig) string {
	fmt.Printf("Generate proto for swift - version %s\n", cfg.PackageVersion)
	dist := protoc.Swift(cfg)
	files := make([]string, 0)
	service := cfg.Project.Name
	namespace := strings.ReplaceAll(cfg.Project.Namespace, "/", "_")
	filepath.Walk(dist, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".swift") {
			files = append(files, path)
		}
		return nil
	})
	for _, f := range files {
		dir, filename := path.Split(f)
		newFilename := fmt.Sprintf("%s_%s_%s", namespace, service, filename)
		newpath := path.Join(dir, newFilename)
		os.Rename(f, newpath)
	}

	return dist
}
