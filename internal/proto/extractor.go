package googleapis

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"
)

//go:embed googleapis.zip
var googleApisZip []byte

func ExtractGoogleApisZip(extractDir string) error {
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return err
	}

	reader := bytes.NewReader(googleApisZip)
	zipReader, err := zip.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		extractPath := filepath.Join(extractDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(extractPath, 0755)
		} else {
			dst, err := os.Create(extractPath)
			if err != nil {
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
