package migrate

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"mykit/internal/config"
	"mykit/internal/constant"
	"mykit/internal/emitter/common"
	"mykit/internal/metadata"
	"mykit/internal/util/gocmd"
	osutil "mykit/internal/util/os"
)

var (
	_oldGoPackage string
	_newGoPackage string

	_oldServiceName string
	_newServiceName string
)

func Migrate(source, newServiceName string, cfg *config.GenerateConfig) {
	_newServiceName = newServiceName

	migrateMyKit(source)

	err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("walk file failed", path, err)
			os.Exit(1)
		}
		if info.IsDir() ||
			info.Name() == "nightkit.yaml" ||
			info.Name() == "go.sum" ||
			info.Name() == "go.mod" ||
			strings.Contains(path, "/.git/") ||
			strings.Contains(path, "/vendor/") ||
			strings.HasSuffix(info.Name(), ".swagger.json") {
			return nil
		}
		for _, f := range _myKitFiles {
			if strings.HasSuffix(path, f.Dest) {
				return nil
			}
		}

		fmt.Println(strings.TrimPrefix(strings.TrimPrefix(path, source), "/"))

		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("read file failed", path, err)
			os.Exit(1)
		}
		newContent := string(content)

		destFileName := info.Name()
		extension := filepath.Ext(destFileName)
		switch extension {
		case ".go":
			newContent = replaceGoImports(newContent)
			if strings.Contains(destFileName, ".pb.") {
				destFileName = renameFileName(destFileName)
			}
		case ".proto":
			newContent = replaceProtoImport(newContent)
			destFileName = renameFileName(destFileName)
		}

		destDir := filepath.Dir(strings.ReplaceAll(path, source, metadata.Dir))
		if err := osutil.CreateDirIfNotExists(destDir, 0755); err != nil {
			fmt.Println("create dir failed", destDir, err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(filepath.Join(destDir, destFileName), []byte(newContent), 0644)
		if err != nil {
			fmt.Println("write file failed", path, err)
			os.Exit(1)
		}

		return nil
	})
	if err != nil {
		fmt.Println("migrate failed", source, err)
	}

	// gocmd.Init(getPackage())
	// gocmd.Get([]string{
	// 	"github.com/ChisTrun/myid@v1.2.1-myid-666fe14c",
	// 	"github.com/ChisTrun/mywallet@v1.0.0-mywallet-cfe7cf0a",
	// 	"github.com/ChisTrun/potter@v1.0.0-potter-c5a6f005",
	// })
	gocmd.Vendor(cfg)
}

func replaceGoImports(fileContent string) string {
	fileContent = strings.ReplaceAll(fileContent, "*api.TCPSocket", "*carbon.TCPSocket")

	if strings.Contains(fileContent, "*carbon.Kafka") {
		fileContent = strings.ReplaceAll(fileContent, "*carbon.Kafka", "*kafkaapi.Kafka")
		if !strings.Contains(fileContent, "*carbon.") {
			if strings.Contains(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"") {
				fileContent = strings.ReplaceAll(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "kafkaapi \"github.com/ChisTrun/kafka/api\"")
			} else {
				fileContent = strings.ReplaceAll(fileContent, "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "kafkaapi \"github.com/ChisTrun/kafka/api\"")
			}
		} else {
			if strings.Contains(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"") {
				fileContent = strings.ReplaceAll(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"\nkafkaapi \"github.com/ChisTrun/kafka/api\"")
			} else {
				fileContent = strings.ReplaceAll(fileContent, "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"\nkafkaapi \"github.com/ChisTrun/kafka/api\"")
			}
		}
	}

	if strings.Contains(fileContent, "*carbon.Redis") {
		fileContent = strings.ReplaceAll(fileContent, "*carbon.Redis", "*redisapi.Redis")
		if !strings.Contains(fileContent, "*carbon.") {
			if strings.Contains(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"") {
				fileContent = strings.ReplaceAll(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "redisapi \"github.com/ChisTrun/redis/api\"")
			} else {
				fileContent = strings.ReplaceAll(fileContent, "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "redisapi \"github.com/ChisTrun/redis/api\"")
			}
		} else {
			if strings.Contains(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"") {
				fileContent = strings.ReplaceAll(fileContent, "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "carbon \"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"\nredisapi \"github.com/ChisTrun/redis/api\"")
			} else {
				fileContent = strings.ReplaceAll(fileContent, "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"", "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon\"\nredisapi \"github.com/ChisTrun/redis/api\"")
			}
		}
	}

	if strings.Contains(fileContent, "nightkitgrpc \"gitlab.com/inspirelab/greyhole/night-kit/pkg/grpc\"") {
		fileContent = strings.ReplaceAll(fileContent, "nightkitgrpc \"gitlab.com/inspirelab/greyhole/night-kit/pkg/grpc\"", "mykitgrpc \"github.com/ChisTrun/grpc/pkg/client\"")
	} else {
		fileContent = strings.ReplaceAll(fileContent, "mykit/pkg/grpc", "grpc \"github.com/ChisTrun/grpc/pkg/client\"")
	}

	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole/night-kit/pkg/config", "github.com/ChisTrun/carbon/pkg/config")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole/night-kit/pkg/carbon", "github.com/ChisTrun/carbon/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole/night-kit/pkg/redis", "github.com/ChisTrun/redis/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole/night-kit/pkg/kafka", "github.com/ChisTrun/kafka/pkg")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole/night-kit/pkg/logging", "github.com/ChisTrun/logger/pkg/logging")
	fileContent = strings.ReplaceAll(fileContent, "\"gitlab.com/inspirelab/greyhole/night-kit/pkg/grpc\"", "grpc \"github.com/ChisTrun/grpc/pkg/client\"")

	fileContent = strings.ReplaceAll(fileContent, "nightkit \"gitlab.com/inspirelab/greyhole/night-kit/pkg/api\"", "mykit \"mykit/pkg/api\"")
	fileContent = strings.ReplaceAll(fileContent, "nightent \"gitlab.com/inspirelab/greyhole/night-kit/pkg/ent\"", "dbe \"github.com/ChisTrun/database/pkg/ent\"")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole/night-kit/pkg/ent", "github.com/ChisTrun/database/pkg/ent")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole", "github.com/ChisTrun")

	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/conveyor/pkg/api", "github.com/ChisTrun/conveyor/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/bentau/pkg/client/v1", "github.com/ChisTrun/bentau/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/playah/pkg/client/v1", "github.com/ChisTrun/playah/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/playah/pkg/client/v2", "github.com/ChisTrun/playah/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/playah/pkg/client/v3", "github.com/ChisTrun/playah/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/playah/pkg/api/v1", "github.com/ChisTrun/playah/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/playah/pkg/api", "github.com/ChisTrun/playah/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/whatsapp/pkg/api", "github.com/ChisTrun/whatsapp/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/whatsapp/pkg/app/v1", "github.com/ChisTrun/whatsapp/pkg/app")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/whatsapp/pkg/app/v2", "github.com/ChisTrun/whatsapp/pkg/app")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/whatsapp/pkg/app/v3", "github.com/ChisTrun/whatsapp/pkg/app")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/bentau/pkg/api", "github.com/ChisTrun/bentau/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/mywallet/pkg/client/steward/v1", "github.com/ChisTrun/mywallet/pkg/client/steward")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/mywallet/pkg/api", "github.com/ChisTrun/mywallet/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/myaccount/pkg/api", "github.com/ChisTrun/myaccount/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/slot/pkg/api", "github.com/ChisTrun/slot/api")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend/tinder/pkg/api", "github.com/ChisTrun/tinder/api")
	fileContent = strings.ReplaceAll(fileContent, "\"gitlab.com/inspirelab/gameloot/monorepo/backend/shared/random\"", "")

	fileContent = strings.ReplaceAll(fileContent, _oldGoPackage, _newGoPackage)
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend", "github.com/ChisTrun")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/backend", "github.com/ChisTrun")

	fileContent = strings.ReplaceAll(fileContent, "github.com/ChisTrun/nats/pkg/pusher/v1", "github.com/ChisTrun/nats/pkg/pusher")
	fileContent = strings.ReplaceAll(fileContent, "mykit/pkg/config", "github.com/ChisTrun/carbon/pkg/config")
	fileContent = strings.ReplaceAll(fileContent, "mykit/pkg/carbon", "github.com/ChisTrun/carbon/api")
	fileContent = strings.ReplaceAll(fileContent, "mykit/pkg/redis", "github.com/ChisTrun/redis/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "mykit/pkg/kafka", "github.com/ChisTrun/kafka/pkg")
	fileContent = strings.ReplaceAll(fileContent, "mykit/pkg/logging", "github.com/ChisTrun/logger/pkg/logging")
	fileContent = strings.ReplaceAll(fileContent, "github.com/ChisTrun/bentau/pkg/client/v1", "github.com/ChisTrun/bentau/pkg/client")
	fileContent = strings.ReplaceAll(fileContent, "github.com/ChisTrun/mywallet/pkg/client/steward/v1", "github.com/ChisTrun/mywallet/pkg/client/steward")

	fileContent = strings.ReplaceAll(fileContent, "*carbon.DatabaseV2", "*carbon.Database")
	fileContent = strings.ReplaceAll(fileContent, "nightkit", "mykit")
	fileContent = strings.ReplaceAll(fileContent, "nightent", "dbe")
	fileContent = strings.ReplaceAll(fileContent, ".OpenV2(", ".Open(")
	fileContent = strings.ReplaceAll(fileContent, "night-kit", "mykit")

	fileContent = strings.ReplaceAll(fileContent, "*"+_oldServiceName+".", "*"+_newServiceName+".")
	fileContent = strings.ReplaceAll(fileContent, "&"+_oldServiceName+".", "&"+_newServiceName+".")
	fileContent = strings.ReplaceAll(fileContent, _oldServiceName+" \"", _newServiceName+" \"")
	if strings.Contains(fileContent, " NewServer") {
		fileContent = strings.ReplaceAll(fileContent, _oldServiceName+".", _newServiceName+".")
	}

	fileContent = strings.ReplaceAll(fileContent, "github.com/ChisTrun/rpc/pkg/api", "github.com/ChisTrun/grpc/pkg/api")
	fileContent = strings.ReplaceAll(fileContent, "github.com/ChisTrun/rpc/pkg/error", "github.com/ChisTrun/grpc/pkg/error")
	fileContent = strings.ReplaceAll(fileContent, "github.com/ChisTrun/rpc/pkg/status", "github.com/ChisTrun/grpc/pkg/status")

	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/dceu/monorepo/backend", "github.com/ChisTrun/dceu/backend")

	return fileContent
}

func replaceProtoImport(fileContent string) string {
	fileContent = strings.ReplaceAll(fileContent, "night-kit", "carbon")
	fileContent = strings.ReplaceAll(fileContent, _oldGoPackage, _newGoPackage)
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/monorepo/backend", "github.com/ChisTrun")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/gameloot/backend", "github.com/ChisTrun")
	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/greyhole", "github.com/ChisTrun")

	if strings.Contains(fileContent, fmt.Sprintf("import \"%s/api/%s_import.proto\";", _oldServiceName, _oldServiceName)) {
		fileContent = strings.ReplaceAll(fileContent, fmt.Sprintf("import \"%s/api/%s_import.proto\";", _oldServiceName, _oldServiceName), "import \"mywallet/api/v1/mywallet_change.proto\";\nimport \"playah/api/playah.proto\";")
		fileContent = strings.ReplaceAll(fileContent, " Change ", " gameloot.mywallet.v1.Change ")
		fileContent = strings.ReplaceAll(fileContent, " Playah ", " gameloot.playah.Playah ")
		fileContent = strings.ReplaceAll(fileContent, " Change>", " gameloot.mywallet.v1.Change>")
		fileContent = strings.ReplaceAll(fileContent, " Playah>", " gameloot.playah.Playah>")
	}

	fileContent = strings.ReplaceAll(fileContent, fmt.Sprintf("import \"%s/api/", _oldServiceName), fmt.Sprintf("import \"%s/api/", strings.Split(_newGoPackage, "/")[len(strings.Split(_newGoPackage, "/"))-1]))
	fileContent = strings.ReplaceAll(fileContent, _oldServiceName, _newServiceName)
	fileContent = strings.ReplaceAll(fileContent, " carbon.NATS", " greyhole.carbon.NATS")
	fileContent = strings.ReplaceAll(fileContent, " carbon.S3", " greyhole.carbon.S3")
	if strings.Contains(fileContent, "carbon.Logger") {
		fileContent = strings.ReplaceAll(fileContent, "import \"validate/validate.proto\";", "import \"validate/validate.proto\";\nimport \"logger/api/logger.proto\";")
		fileContent = strings.ReplaceAll(fileContent, "carbon.Logger", "greyhole.logger.Logger")
	}

	if strings.Contains(fileContent, "carbon.DatabaseV2") {
		fileContent = strings.ReplaceAll(fileContent, "import \"validate/validate.proto\";", "import \"validate/validate.proto\";\nimport \"database/api/database.proto\";")
		fileContent = strings.ReplaceAll(fileContent, "carbon.DatabaseV2", "greyhole.database.Database")
	}

	if strings.Contains(fileContent, "carbon.Redis") {
		fileContent = strings.ReplaceAll(fileContent, "import \"validate/validate.proto\";", "import \"validate/validate.proto\";\nimport \"redis/api/redis.proto\";")
		fileContent = strings.ReplaceAll(fileContent, "carbon.Redis", "greyhole.redis.Redis")
	}

	if strings.Contains(fileContent, "carbon.Kafka") {
		fileContent = strings.ReplaceAll(fileContent, "import \"validate/validate.proto\";", "import \"validate/validate.proto\";\nimport \"kafka/api/kafka.proto\";")
		fileContent = strings.ReplaceAll(fileContent, "carbon.Kafka", "greyhole.kafka.Kafka")
	}

	fileContent = strings.ReplaceAll(fileContent, "carbon.Listener", "greyhole.carbon.Listener")
	fileContent = strings.ReplaceAll(fileContent, "carbon.TCPSocket", "greyhole.carbon.TCPSocket")

	fileContent = strings.ReplaceAll(fileContent, "gitlab.com/inspirelab/dceu/monorepo/backend", "github.com/ChisTrun/dceu/backend")

	return fileContent
}

func renameFileName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, _oldServiceName, _newServiceName)

	if strings.Contains(fileName, "_"+_newServiceName+".") || strings.Contains(fileName, "_"+_newServiceName+"_") {
		fileName = strings.ReplaceAll(fileName, "_"+_newServiceName, "")
		fileName = _newServiceName + "_" + fileName
	}

	return fileName
}

type OldConfig struct {
	Project struct {
		Name string `yaml:"name"`
	}
	Package string `yaml:"package"`
	Gen     struct {
		Ent struct {
			Enable bool `yaml:"enable"`
		} `yaml:"ent"`
		Client struct {
			Enable bool `yaml:"enable"`
		} `yaml:"client"`
		Proto struct {
			Go []string `yaml:"go"`
			Js []string `yaml:"js"`
		} `yaml:"proto"`
	}
}

var _myKitFiles = []struct {
	Src             string
	Dest            string
	SkipForMonorepo bool
}{
	{
		Src:             "goservice/build/Dockerfile.tmpl",
		Dest:            "build/Dockerfile",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/build/local.Dockerfile.tmpl",
		Dest:            "build/local.Dockerfile",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/mykit.yaml.tmpl",
		Dest:            "mykit.yaml",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/Makefile.tmpl",
		Dest:            "Makefile",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/gitignore.tmpl",
		Dest:            ".gitignore",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/dockerignore.tmpl",
		Dest:            ".dockerignore",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/ci.json.tmpl",
		Dest:            "ci.json",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/gitlab-ci.yml.tmpl",
		Dest:            ".gitlab-ci.yml",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/api/code.proto.tmpl",
		Dest:            "api/xtype_code.proto",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/go.mod.tmpl",
		Dest:            "go.mod",
		SkipForMonorepo: true,
	},
}

func migrateMyKit(source string) {
	nightkitPath := filepath.Join(source, "nightkit.yaml")
	data, err := ioutil.ReadFile(nightkitPath)
	if err != nil {
		fmt.Println("read file failed", nightkitPath, err)
		os.Exit(1)
	}

	oldConfig := OldConfig{}
	err = yaml.Unmarshal(data, &oldConfig)
	if err != nil {
		fmt.Println("unmarshal file failed", nightkitPath, err)
		os.Exit(1)
	}

	_oldServiceName = oldConfig.Project.Name
	if len(_newServiceName) == 0 {
		_newServiceName = _oldServiceName
	}
	_oldGoPackage = oldConfig.Package
	if len(_oldGoPackage) == 0 {
		_oldGoPackage = strings.ReplaceAll(source, constant.GoPath, "")
	}
	_newGoPackage = getPackage()

	var (
		generateGo []string
		generateJs []string
		monorepo   bool
	)
	var goHasCode bool
	for _, fileName := range oldConfig.Gen.Proto.Go {
		newFileName := renameFileName(fileName)
		generateGo = append(generateGo, newFileName)
		if strings.HasSuffix(newFileName, "_code.proto") {
			goHasCode = true
		}
	}
	if !goHasCode {
		generateGo = append(generateGo, _newServiceName+"_code.proto")
	}

	var jsHasCode bool
	for _, fileName := range oldConfig.Gen.Proto.Js {
		newFileName := renameFileName(fileName)
		generateJs = append(generateJs, newFileName)
		if strings.HasSuffix(newFileName, "_code.proto") {
			jsHasCode = true
		}
	}
	if !jsHasCode {
		generateJs = append(generateJs, _newServiceName+"_code.proto")
	}

	//if _, err := os.Stat(filepath.Join(source, "go.mod")); err != nil {
	//	monorepo = true
	//}

	for _, f := range _myKitFiles {
		if monorepo && f.SkipForMonorepo {
			continue
		}

		dest := strings.ReplaceAll(f.Dest, "xtype", _newServiceName)
		common.Render(f.Src, filepath.Join(metadata.Dir, dest),
			map[string]interface{}{
				"ProjectName": _newServiceName,
				"Package":     _newGoPackage,
				"Monorepo":    monorepo,
				"Ent":         oldConfig.Gen.Ent.Enable,
				"Client":      oldConfig.Gen.Client.Enable,
				"GenerateGo":  generateGo,
				"GenerateJs":  generateJs,
			})
	}
}

func getPackage() string {
	if !strings.HasPrefix(metadata.Dir, constant.GoPath) {
		return _newServiceName
	}
	pkg, err := filepath.Rel(constant.GoPath, metadata.Dir)
	if err != nil {
		fmt.Println("get default package failed", err)
		return _newServiceName
	}

	return pkg
}
