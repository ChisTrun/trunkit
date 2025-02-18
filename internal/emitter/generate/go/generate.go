package _go

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"gopkg.in/yaml.v3"

	"trunkit/internal/config"
	"trunkit/internal/constant"
	"trunkit/internal/emitter/common"
	"trunkit/internal/metadata"
	"trunkit/internal/protoc"
	"trunkit/internal/template"
	osutil "trunkit/internal/util/os"
)

const (
	_clientTemplate = "goservice/pkg/client/xservice.go.tmpl"
	_clientDest     = "pkg/client/xServiceName.go"
)

type protoService struct {
	GoPackage     string
	GoPackageName string
	GoFilesPath   string
	Services      []string
}

type enumExtension struct {
	Name          string
	AllowEnumZero bool
	GoPackageName string
	GoFilesPath   string
}

func Generate(cfg *config.GenerateConfig) {
	cfgStr, _ := yaml.Marshal(cfg)
	fmt.Println("Config:")
	fmt.Println(string(cfgStr))
	protoc.Go(cfg)
	fmt.Println("Generate GRPC apis for Go:")
	var (
		protoServices  []*protoService
		enumExtensions []*enumExtension
	)

	for _, fileName := range cfg.Generate.Proto.Go {
		filePath := filepath.Join(metadata.Dir, "api", fileName)
		definition, err := protoc.ParseProto(filePath)
		if err != nil {
			fmt.Println("parse file failed", filePath, err)
			os.Exit(1)
		}

		var (
			goPackage     string
			goFilesPath   string
			goPackageName string
			services      []string
			protoRPCs     = map[string][]*proto.RPC{}
		)
		proto.Walk(definition,
			proto.WithService(func(service *proto.Service) {
				services = append(services, service.Name)
				if _, found := protoRPCs[service.Name]; !found {
					protoRPCs[service.Name] = []*proto.RPC{}
				}
			}),
			proto.WithRPC(func(rpc *proto.RPC) {
				serviceName := rpc.Parent.(*proto.Service).Name
				protoRPCs[serviceName] = append(protoRPCs[serviceName], rpc)
			}),
			proto.WithOption(func(option *proto.Option) {
				if option.Name == "go_package" {
					goPackageValues := strings.Split(option.Constant.Source, ";")
					goPackageName = goPackageValues[1]
					goFilesPath = osutil.MergePath(metadata.Dir, goPackageValues[0])
					goPackage = filepath.Join(cfg.Project.GoPackage, strings.TrimPrefix(goFilesPath, metadata.Dir))
				}
			}),
			proto.WithEnum(func(enum *proto.Enum) {
				enumConfig, found := cfg.Generate.Proto.GoAdditionalEnum[enum.Name]
				if !found {
					return
				}

				enumExtensions = append(enumExtensions, &enumExtension{
					Name:          enum.Name,
					AllowEnumZero: enumConfig.AllowEnumZero,
					GoPackageName: goPackageName,
					GoFilesPath:   goFilesPath,
				})
			}),
		)

		if len(protoRPCs) > 0 {
			fmt.Println("-", fileName)
			generateRPCServices(protoRPCs, goPackage, cfg)

			var existed bool
			for _, s := range protoServices {
				if s.GoPackage == goPackage {
					s.Services = append(s.Services, services...)
					existed = true
					break
				}
			}
			if !existed {
				protoServices = append(protoServices, &protoService{
					GoPackage:     goPackage,
					GoPackageName: goPackageName,
					Services:      services,
					GoFilesPath:   goFilesPath,
				})
			}
		}
	}
	fmt.Println()
	generateCommand(cfg)
	generateEnumExtensions(enumExtensions)
	generateServer(protoServices, cfg)
	generateClients(protoServices, cfg)
	generateDockerfile(cfg)

	// gocmd.Vendor(cfg)~

	generateEnt(cfg)
	generateClientMocks(protoServices)

	// gocmd.Vendor(cfg)
}

const (
	_serviceTemplate = "goservice/internal/server/service/grpc_xtype.go.tmpl"
	_serviceDest     = "internal/server/xServiceName/grpc_xServiceName.go"
	_grpcDefaultPath = "internal/server"
)

func generateRPCServices(serviceRPCs map[string][]*proto.RPC, goPackage string, cfg *config.GenerateConfig) {
	generateCfg := cfg.Generate
	if generateCfg.GrpcServer.Enable {
		fmt.Println("Generate grpc server:")

		for serviceName, protoRPCs := range serviceRPCs {
			fmt.Println(" +", serviceName)

			serviceDest := strings.ReplaceAll(_serviceDest, "xServiceName", strings.ToLower(serviceName))
			serviceDest = strings.ReplaceAll(serviceDest, _grpcDefaultPath, cfg.Generate.GrpcServer.Path)
			common.Render(
				_serviceTemplate,
				serviceDest,
				map[string]interface{}{
					"ProjectName": cfg.Project.Name,
					"Package":     goPackage,
					"ServiceName": serviceName,
				},
			)

			for _, rpc := range protoRPCs {
				fmt.Println("  ", rpc.Name, "-", generateRPC(rpc, goPackage, cfg))
			}
		}
	}
}

const (
	_methodTemplate = "goservice/internal/server/service/grpc_xmethod.go.tmpl"
	_methodDest     = "internal/server/xServiceName/grpc_xServiceName_xServiceMethod.go"
)

func generateRPC(rpc *proto.RPC, goPackage string, cfg *config.GenerateConfig) string {

	imports := make(map[string]*goType)

	requestType, err := getProtoType(cfg.Project.Name, rpc.RequestType, cfg)
	if err != nil {
		fmt.Println()
		fmt.Println("generate rpc failed.", err)
		os.Exit(1)
	}

	returnsType, err := getProtoType(cfg.Project.Name, rpc.ReturnsType, cfg)
	if err != nil {
		fmt.Println()
		fmt.Println("generate rpc failed.", err)
		os.Exit(1)
	}

	if len(requestType.ImportPath) > 0 {
		imports[requestType.ImportPath] = requestType
	}
	if len(returnsType.ImportPath) > 0 {
		imports[returnsType.ImportPath] = returnsType
	}

	protoServiceName := rpc.Parent.(*proto.Service).Name
	methodDest := strings.NewReplacer(
		"xServiceName", strings.ToLower(protoServiceName),
		"xServiceMethod", strings.ToLower(rpc.Name)).
		Replace(_methodDest)
	methodDest = strings.ReplaceAll(methodDest, _grpcDefaultPath, cfg.Generate.GrpcServer.Path)
	common.Render(
		_methodTemplate,
		methodDest,
		map[string]interface{}{
			"Package":        goPackage,
			"ServiceName":    protoServiceName,
			"Method":         rpc.Name,
			"RequestType":    requestType.Value,
			"ReturnsType":    returnsType.Value,
			"Imports":        imports,
			"HasImportProto": len(requestType.ImportPath) == 0 || len(returnsType.ImportPath) == 0,
		},
	)

	return methodDest
}

const (
	_serveTemplate  = "goservice/internal/server/serve.go.tmpl"
	_serveDest      = "internal/server/serve.go"
	_serverTemplate = "goservice/internal/server/z_server.go.tmpl"
	_serverDest     = "internal/server/z_server.go"
	_cmdTemplate    = "goservice/cmd/main.go.tmpl"
	_cmdDest        = "cmd/main.go"
)

func generateServer(protoServices []*protoService, cfg *config.GenerateConfig) {

	generateCfg := cfg.Generate
	if generateCfg.Server.Enable {
		serveDest := path.Join(generateCfg.Server.Path, "serve.go")
		fmt.Printf("Generate server: %s\n", serveDest)
		common.Render(
			_serveTemplate,
			serveDest,
			map[string]interface{}{
				"MyKitBase":      constant.MyKitBase,
				"ProjectName":    cfg.Project.Name,
				"ServiceName":    cfg.Project.Name,
				"Package":        cfg.Project.GoPackage,
				"GrpcServerPath": generateCfg.GrpcServer.Path,
				"ProtoServices":  protoServices,
				"GenOpts":        cfg.Generate,
				"GoConfigPath":   cfg.Project.GoConfigPath,
			})

		serverDest := path.Join(generateCfg.Server.Path, "z_server.go")
		fmt.Printf("Generate server: %s\n", serverDest)
		common.Render(
			_serverTemplate,
			serverDest,
			map[string]interface{}{
				"Version":      metadata.MyKitVersion,
				"MyKitBase":    constant.MyKitBase,
				"Package":      cfg.Project.GoPackage,
				"GenOpts":      cfg.Generate,
				"GoConfigPath": cfg.Project.GoConfigPath,
			},
			common.Overwrite())
	}

}

func generateCommand(cfg *config.GenerateConfig) {
	generateCfg := cfg.Generate
	if generateCfg.Command.Enable {
		fmt.Println("Generate command")
		common.Render(
			_cmdTemplate,
			cfg.Generate.Command.Path,
			map[string]interface{}{
				"Package": cfg.Project.GoPackage,
			})
	}
}

func generateClients(protoServices []*protoService, cfg *config.GenerateConfig) {
	if len(protoServices) == 0 {
		return
	}

	if cfg.Generate.Client.Enable {
		if err := osutil.CreateDirIfNotExists(filepath.Join(metadata.Dir, "pkg/client"), 0755); err != nil {
			os.Exit(1)
		}

		for _, service := range protoServices {
			for _, serviceName := range service.Services {
				common.Render(
					_clientTemplate,
					strings.ReplaceAll(_clientDest, "xServiceName", strings.ToLower(serviceName)),
					map[string]interface{}{
						"Version":     metadata.MyKitVersion,
						"ServiceName": serviceName,
						"Package":     service.GoPackage,
					},
					common.Overwrite(),
				)
			}
		}
	}
}

const (
	_clientDocTemplate = "goservice/pkg/api/doc.go.tmpl"
	_clientDocDest     = "doc.go"
)

func generateClientMocks(protoServices []*protoService) {
	fmt.Println("Generate mocks:")
	for _, protoService := range protoServices {
		fmt.Println("-", protoService.GoFilesPath)
		for _, serviceName := range protoService.Services {
			fmt.Println(" +", serviceName)
		}

		common.Render(
			_clientDocTemplate,
			filepath.Join(protoService.GoFilesPath, _clientDocDest),
			map[string]interface{}{
				"PackageName": protoService.GoPackageName,
				"Services":    protoService.Services,
			},
			common.Overwrite(),
		)

		_, err := osutil.Exec(
			[]string{
				fmt.Sprintf("cd %s", protoService.GoFilesPath),
				"go generate .",
			},
		)
		if err != nil {
			fmt.Println("go generate mock docs failed", err)
		}
	}
}

const (
	_dockerfileTemplate      = "goservice/build/Dockerfile.tmpl"
	_dockerfileDest          = "build/Dockerfile"
	_localDockerfileTemplate = "goservice/build/local.Dockerfile.tmpl"
	_localDockerfileDest     = "build/local.Dockerfile"
)

func generateDockerfile(cfg *config.GenerateConfig) {
	if cfg.Generate.DockerFile.Enable {
		fmt.Println("Generate dockerfile")
		common.Render(
			_dockerfileTemplate,
			_dockerfileDest,
			map[string]interface{}{
				"Monorepo": cfg.Project.Monorepo,
			},
		)
		common.Render(
			_localDockerfileTemplate,
			_localDockerfileDest,
			map[string]interface{}{
				"Monorepo": cfg.Project.Monorepo,
			},
		)
	}
}

const (
	_pkgDocTemplate = "goservice/pkg/doc.go.tmpl"
	_pkgDocDest     = "pkg/doc.go"
)

func generateEnt(cfg *config.GenerateConfig) {
	if cfg.Generate.Ent.Enable {
		common.Render(
			_pkgDocTemplate,
			_pkgDocDest,
			map[string]interface{}{},
		)

		_, err := osutil.Exec(
			[]string{
				fmt.Sprintf("cd %s", filepath.Join(metadata.Dir, "pkg")),
				"go generate .",
			},
		)
		if err != nil {
			fmt.Println("go generate ent failed", err)
		}
	}
}

type goType struct {
	Value       string
	ImportPath  string
	ImportAlias string
}

func getProtoType(projectName string, typeStr string, cfg *config.GenerateConfig) (*goType, error) {
	switch {
	case typeStr == "google.protobuf.Empty":
		return &goType{
			Value:      "emptypb.Empty",
			ImportPath: "google.golang.org/protobuf/types/known/emptypb",
		}, nil
	case typeStr == "google.protobuf.Timestamp":
		return &goType{
			Value:      "time.Time",
			ImportPath: "google.golang.org/protobuf/types/known/timestamppb",
		}, nil
	case typeStr == "google.protobuf.Any":
		return &goType{
			Value:      "any.Any",
			ImportPath: "google.golang.org/protobuf/types/known/anypb",
		}, nil
	case typeStr == "google.protobuf.Struct":
		return &goType{
			Value:      "_struct.Struct",
			ImportPath: "google.golang.org/protobuf/types/known/structpb",
		}, nil
	case strings.Contains(typeStr, "."):
		typeElements := strings.Split(typeStr, ".")
		messageName := typeElements[len(typeElements)-1]
		importAlias := typeElements[len(typeElements)-2]

		for _, i := range cfg.Generate.Proto.Imports {
			for _, requestType := range i.Types {
				if requestType == typeStr {
					return &goType{
						Value:       fmt.Sprintf("%s.%s", importAlias, messageName),
						ImportAlias: importAlias,
						ImportPath:  i.GoPackage,
					}, nil
				}
			}
			// for _, requestType := range i.Requests {
			// 	if requestType == typeStr {
			// 		return &goType{
			// 			Value:       fmt.Sprintf("%s.%s", importAlias, messageName),
			// 			ImportAlias: importAlias,
			// 			ImportPath:  i.GoPackage,
			// 		}, nil
			// 	}
			// }

			// for _, returnType := range i.Returns {
			// 	if returnType == typeStr {
			// 		return &goType{
			// 			Value:       fmt.Sprintf("%s.%s", importAlias, messageName),
			// 			ImportAlias: importAlias,
			// 			ImportPath:  i.GoPackage,
			// 		}, nil
			// 	}
			// }
		}
	default:
		return &goType{
			Value: fmt.Sprintf("%s.%s", projectName, typeStr),
		}, nil
	}

	return nil, fmt.Errorf("import %s not found in", typeStr)
}

const (
	_enumTemplate = "goservice/api/xenum.go.tmpl"
	_enumDest     = "/xenum.go"
)

func generateEnumExtensions(enumExtensions []*enumExtension) {
	for _, enumExtension := range enumExtensions {
		enumExtensionDest := strings.ReplaceAll(_enumDest, "xenum", template.CamelToSnake(enumExtension.Name))
		common.Render(
			_enumTemplate,
			filepath.Join(enumExtension.GoFilesPath, enumExtensionDest),
			map[string]interface{}{
				"Version":       metadata.MyKitVersion,
				"PackageName":   enumExtension.GoPackageName,
				"Enum":          enumExtension.Name,
				"AllowEnumZero": enumExtension.AllowEnumZero,
			},
			common.Overwrite(),
		)
	}
}
