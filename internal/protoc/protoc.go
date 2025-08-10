package protoc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"

	"github.com/ChisTrun/trunkit/internal/config"
	"github.com/ChisTrun/trunkit/internal/metadata"
	googleapis "github.com/ChisTrun/trunkit/internal/proto"
	"github.com/ChisTrun/trunkit/internal/util/gocmd"
	osutil "github.com/ChisTrun/trunkit/internal/util/os"
)

func ParseProto(filePath string) (*proto.Proto, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file failed", filePath, err)
		return nil, err
	}

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Println("parse file failed", filePath, err)
		return nil, err
	}

	return definition, nil
}

func GetTypeName(v proto.Visitee) string {
	switch v.(type) {
	case *proto.Message:
		return v.(*proto.Message).Name
	case *proto.Enum:
		return v.(*proto.Enum).Name
	case *proto.Package:
		return v.(*proto.Package).Name
	case *proto.Service:
		return v.(*proto.Service).Name
	case *proto.RPC:
		return v.(*proto.RPC).Name
	case *proto.Option:
		return v.(*proto.Option).Name
	case *proto.Proto:
		pkg := ""
		proto.Walk(v.(*proto.Proto), proto.WithPackage(func(p *proto.Package) {
			pkg = p.Name
		}))
		return pkg

	default:
		return ""
	}
}

func GetParent(v proto.Visitee) proto.Visitee {
	switch v.(type) {
	case *proto.Message:
		return v.(*proto.Message).Parent
	case *proto.Enum:
		return v.(*proto.Enum).Parent
	case *proto.Package:
		return v.(*proto.Package).Parent
	case *proto.Service:
		return v.(*proto.Service).Parent
	case *proto.RPC:
		return v.(*proto.RPC).Parent
	case *proto.Option:
		return v.(*proto.Option).Parent
	case *proto.Oneof:
		return v.(*proto.Oneof).Parent
	case *proto.MapField:
		return v.(*proto.MapField).Parent
	default:
		return nil
	}
}

func GetFullTypeName(m proto.Visitee) string {
	name := GetTypeName(m)
	p := GetParent(m)
	for p != nil {
		pName := GetTypeName(p)
		if pName != "" {
			name = pName + "." + name
		}
		p = GetParent(p)
	}
	return name
}

func GetAllTypes(m *proto.Proto) []string {
	var types []string
	proto.Walk(m, proto.WithMessage(func(m *proto.Message) {
		types = append(types, GetFullTypeName(m))
	}), proto.WithEnum(func(e *proto.Enum) {
		types = append(types, GetFullTypeName(e))
	}))

	proto.Walk(m, proto.WithEnum(func(e *proto.Enum) {
		types = append(types, GetFullTypeName(e))
	}))

	return types
}

func executeProtoc(context string, protoFiles []string, outPaths []string, cfg *config.GenerateConfig) {
	var args []string
	for _, i := range getImports(protoFiles, cfg) {
		args = append(args, "-I", i)
	}
	args = append(args, outPaths...)

	for _, fileName := range protoFiles {
		filePath := fmt.Sprintf("%s/api/%s", metadata.Dir, fileName)
		if _, err := os.Stat(filePath); err != nil {
			fmt.Println("proto file", filePath, err)
			os.Exit(1)
		}
		args = append(args, filePath)
	}

	fmt.Println("protoc")
	for i := 0; i < len(args); i++ {
		if args[i] == "-I" {
			fmt.Printf("%s %s\n", args[i], args[i+1])
			i++
		} else {
			fmt.Println(args[i])
		}
	}
	fmt.Println()

	_, err := osutil.Exec(
		[]string{
			fmt.Sprintf("cd %s", context),
			strings.Join(append([]string{"protoc"}, args...), " "),
		},
	)
	if err != nil {
		fmt.Println("execute proto c failed", err)
		os.Exit(1)
	}

	gocmd.Vendor(cfg)
}

func readAllTypeInProtos(protoFile string) []string {
	return nil
}

func getImports(protoFiles []string, cfg *config.GenerateConfig) []string {
	var (
		importMap  = map[string]struct{}{}
		imports    []string
		goPackages []string
		// types      map[string][]string
	)

	vendorDir := metadata.Dir
	if cfg.Project.Monorepo {
		vendorDir = filepath.Dir(metadata.Dir)
	}

	defaultProtoImports, defaultGoPackageImports := getDefaultImports(protoFiles)
	for _, i := range defaultProtoImports {
		importPath := filepath.Join(vendorDir, "vendor", i)
		if _, found := importMap[importPath]; !found {
			imports = append(imports, importPath)
			importMap[importPath] = struct{}{}
		}
	}
	goPackages = append(goPackages, defaultGoPackageImports...)

	for _, i := range cfg.Generate.Proto.Imports {
		goPackages = append(goPackages, i.GoPackage)
		importPath := filepath.Join(vendorDir, "vendor", i.Path)
		if _, found := importMap[importPath]; !found {
			imports = append(imports, importPath)
			importMap[importPath] = struct{}{}
		}
	}

	gocmd.Vendor(cfg)
	gocmd.Bootstrap(goPackages, cfg)

	// Resolve proto types of import files

	if cfg.Project.Monorepo {
		vendorDir = filepath.Dir(metadata.Dir)
	}

	for k, i := range cfg.Generate.Proto.Imports {
		filePath := filepath.Join(vendorDir, "vendor", i.Path, k)
		definition, err := ParseProto(filePath)
		if err != nil {
			fmt.Println("parse file failed", filePath, err)
			os.Exit(1)
		}
		types := GetAllTypes(definition)

		i.Types = types
	}

	if _, found := importMap[filepath.Dir(metadata.Dir)]; !found {
		imports = append(imports, filepath.Dir(metadata.Dir))
	}

	if cfg.Generate.GrpcGateway.Enable {
		extractPath := filepath.Join(vendorDir, "vendor")
		importPath := filepath.Join(extractPath, "googleapis")
		err := googleapis.ExtractGoogleApisZip(extractPath)
		if err != nil {
			fmt.Println("extract googleapis failed", err)
			os.Exit(1)
		}

		imports = append(imports, importPath)
	}

	return imports
}

// key is proto file name
// value is the path to proto file
// generated go files should be in the same place with proto file
var _defaultImportPaths = map[string]string{
	"carbon.proto":   "github.com/ChisTrun/carbon/api",
	"database.proto": "github.com/ChisTrun/database/api",
	"logger.proto":   "github.com/ChisTrun/logger/api",
	"redis.proto":    "github.com/ChisTrun/redis/api",
	"kafka.proto":    "github.com/ChisTrun/kafka/api",
}

var _defaultGoImports = []string{
	"github.com/envoyproxy/protoc-gen-validate/validate",
	"github.com/ChisTrun/carbon/api",
	"github.com/ChisTrun/database/api",
	"github.com/ChisTrun/logger/api",
	"github.com/ChisTrun/redis/api",
	"github.com/ChisTrun/kafka/api",
}

func getDefaultImports(protoFiles []string) ([]string, []string) {
	protoImportMap := map[string]struct{}{
		"github.com/ChisTrun":                       {},
		"github.com/envoyproxy/protoc-gen-validate": {},
	}

	goPackageMap := map[string]struct{}{}
	for _, goPackage := range _defaultGoImports {
		goPackageMap[goPackage] = struct{}{}
	}

	for _, fileName := range protoFiles {
		filePath := filepath.Join(metadata.Dir, "api", fileName)
		definition, err := ParseProto(filePath)
		if err != nil {
			fmt.Println("parse file failed", filePath, err)
			os.Exit(1)
		}

		proto.Walk(definition,
			proto.WithImport(func(i *proto.Import) {
				for protoFileName, protoPath := range _defaultImportPaths {
					if !strings.HasSuffix(i.Filename, protoFileName) {
						continue
					}
					goPackageMap[protoPath] = struct{}{}

					protoPath = filepath.Clean(
						strings.TrimSuffix(
							protoPath, filepath.Clean(
								strings.TrimSuffix(i.Filename, protoFileName),
							),
						),
					)
					protoImportMap[protoPath] = struct{}{}
				}
			}),
		)
	}

	var protoImports, goPackages []string
	for protoImport := range protoImportMap {
		protoImports = append(protoImports, protoImport)
	}
	for goPackage := range goPackageMap {
		goPackages = append(goPackages, goPackage)
	}

	return protoImports, goPackages
}

func getGoogleProtoImport() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "googleapis")
}
