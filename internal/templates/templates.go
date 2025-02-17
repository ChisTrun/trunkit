package templates

import (
	"embed"
	"path/filepath"
)

//go:embed *
var f embed.FS

func Load(templateStr string) ([]byte, error) {
	src := filepath.Join(templateStr)

	content, err := f.ReadFile(src)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}
