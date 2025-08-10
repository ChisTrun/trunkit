package protoc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"

	"github.com/ChisTrun/trunkit/internal/config"
	"github.com/ChisTrun/trunkit/internal/metadata"
	osutil "github.com/ChisTrun/trunkit/internal/util/os"
)

func Go(cfg *config.GenerateConfig) {
	protoFiles := cfg.Generate.Proto.Go
	if len(protoFiles) == 0 {
		return
	}

	definitions := []*proto.Proto{}
	for _, name := range protoFiles {
		result, err := ParseProto(filepath.Join(metadata.Dir, "api", name))
		if err != nil {
			fmt.Println("parse proto file failed", name, err)
			os.Exit(1)
		}
		definitions = append(definitions, result)
	}

	var (
		outPath      string
		useCustomTag bool
	)
	for _, def := range definitions {
		proto.Walk(def,
			proto.WithOption(func(option *proto.Option) {
				if option.Name == "go_package" && len(outPath) == 0 {
					goPackage := strings.Split(option.Constant.Source, ";")[0]
					outPath = osutil.TrimPathSuffix(metadata.Dir, goPackage)
				}
			}),
			proto.WithImport(func(i *proto.Import) {
				if i.Filename == "tagger/tagger.proto" {
					useCustomTag = true
				}
			}),
		)

		if len(outPath) > 0 && useCustomTag {
			break
		}
	}

	outPaths := []string{
		fmt.Sprintf("--go_out=%s", outPath),
		fmt.Sprintf("--go-grpc_out=%s", outPath),
		fmt.Sprintf("--validate_out=lang=go:%s", outPath),
	}

	if cfg.Generate.GrpcGateway.Enable {
		outPaths = append(outPaths, fmt.Sprintf("--grpc-gateway_out=%s", outPath))
	}

	fmt.Println("Generate proto files for Go:")
	executeProtoc(".", cfg.Generate.Proto.Go, outPaths, cfg)

	if useCustomTag {
		replaceProtoTagPaths := []string{
			fmt.Sprintf("--gotag_out=outabsdir=\"%s:%s\"", outPath, outPath),
		}

		// this needs to run after the first command has run, otherwise it will fail/produce unexpected results
		fmt.Println("Adding custom tags for generated proto structs:")
		executeProtoc(".", cfg.Generate.Proto.Go, replaceProtoTagPaths, cfg)
	}
}
