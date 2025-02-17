package config

type BaseConfig struct {
	Version string `default:"1" yaml:"version"`
	Extend  string `yaml:"extend"`
}

type ProjectBaseConfig struct {
	Name          string `yaml:"name"`
	Namespace     string `yaml:"namespace"`
	Monorepo      bool   `yaml:"monorepo"`
	GoPackage     string `yaml:"go_package"`
	NpmPackage    string `yaml:"npm_package"`
	NpmRegistry   string `yaml:"npm_registry"`
	MavenRegistry string `yaml:"maven_registry"`
	GoConfigPath  string `default:"pkg/config" yaml:"go_config_path"`
}

type GenerateBaseConfig struct {
	AllowCustomOptions bool `yaml:"allow_custom_options"`
	Profiling          struct {
		Enable bool   `default:"false" yaml:"enable"`
		Port   string `default:"6060" yaml:"port"`
	} `yaml:"profiling"`
	Command struct {
		Enable bool   `default:"true" yaml:"enable"`
		Path   string `default:"cmd/main.go" yaml:"path"`
	} `yaml:"command"`
	Server struct {
		Enable bool   `default:"true" yaml:"enable"`
		Path   string `default:"internal/server" yaml:"path"`
	} `yaml:"server"`
	HttpServer struct {
		Enable bool `default:"false" yaml:"enable"`
	} `yaml:"http_server"`
	GrpcServer struct {
		Enable bool   `default:"true" yaml:"enable"`
		Path   string `default:"internal/server" yaml:"path"`
	} `yaml:"grpc_server"`
	GrpcGateway struct {
		Enable bool `default:"false" yaml:"enable"`
	} `yaml:"grpc_gateway"`
	DockerFile struct {
		Enable bool `yaml:"enable"`
	} `yaml:"dockerfile"`
	Ent struct {
		Enable bool `yaml:"enable"`
	} `yaml:"ent"`
	Client struct {
		Enable bool   `yaml:"enable"`
		Path   string `default:"pkg/client"  yaml:"path"`
	} `yaml:"client"`
	Helm struct {
		Enable bool `yaml:"enable"`
	} `yaml:"helm"`
	GrpcLog struct {
		Enable bool `yaml:"enable"`
	} `yaml:"grpc_log"`
	OpenMetrics struct {
		Enable bool   `yaml:"enable"`
		Path   string `yaml:"path"`
		Port   int    `yaml:"port"`
	} `yaml:"open_metrics"`
}

type ProtoBaseConfig struct {
	GoAdditionalEnum map[string]struct {
		AllowEnumZero bool `yaml:"allow_enum_zero"`
	} `yaml:"go_additional_enum"`
}

type ImportConfig struct {
	Path            string   `yaml:"path"`
	GoPackage       string   `yaml:"go_package"`
	NpmPackage      string   `yaml:"npm_package"`
	NpmRegistry     string   `yaml:"npm_registry"`
	PackageRegistry string   `yaml:"registry"`
	MavenRegistry   string   `yaml:"maven_registry"`
	JavaPackage     string   `yaml:"java_package"`
	Types           []string `yaml:"types"` // auto detect
}
