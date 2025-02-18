package publish

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	config "trunkit/internal/config"
	"trunkit/internal/emitter/generate/swift"
	"trunkit/internal/metadata"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"
)

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(path.Join(basePath, file.Name()))
		if !file.IsDir() {
			dat, err := os.ReadFile(path.Join(basePath, file.Name()))
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func PackageSwift(cfg *config.GenerateConfig) (string, error) {
	dir := swift.Generate(cfg)
	fmt.Println("creating zip archive...")
	archivePath := path.Join(dir, "..", fmt.Sprintf("%s-%s.zip", cfg.Project.Name, cfg.PackageVersion))
	archive, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	addFiles(zipWriter, dir, fmt.Sprintf("%s/", "api"))
	info := &PackageInfo{
		Name:         cfg.Project.Name,
		Namespace:    cfg.Project.Namespace,
		Version:      cfg.PackageVersion,
		MykitVersion: metadata.MyKitVersion,
	}

	versionWriter, err := zipWriter.Create("package-info.yml")
	if err != nil {
		return "", err
	}
	encoder := yaml.NewEncoder(versionWriter)
	encoder.SetIndent(2)
	err = encoder.Encode(info)
	return archivePath, err
}

func UploadSwift(bucket, region, localPath, remotePath string) error {
	cfg, err := awsCfg.LoadDefaultConfig(context.TODO(), awsCfg.WithRegion(region))
	if err != nil {
		log.Printf("init aws config error: %v", err)
		return err
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remotePath),
		Body:   f,
	})

	return err
}
func PublishSwift(bucket, region string, cfg *config.GenerateConfig) error {
	if bucket == "" || region == "" {
		return fmt.Errorf("bucket or region is empty")
	}
	// Generate & package proto into zip files
	zipFile, err := PackageSwift(cfg)
	if err != nil {
		return err
	}

	// Upload to s3
	service := cfg.Project.Name
	namespace := cfg.Project.Namespace
	_, fileName := filepath.Split(zipFile)

	builder := strings.Builder{}
	if namespace != "" {
		builder.WriteString(namespace)
		builder.WriteString("/")
	}
	builder.WriteString(service)
	builder.WriteString("/")
	builder.WriteString(fileName)
	remoteFile := builder.String()

	err = UploadSwift(bucket, region, zipFile, remoteFile)
	return err
}
