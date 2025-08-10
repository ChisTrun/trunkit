package constant

import (
	"go/build"
	"path/filepath"
)

var (
	GoPath    = filepath.Join(build.Default.GOPATH, "src")
	GoBin     = filepath.Join(build.Default.GOPATH, "bin")
	GoInclude = filepath.Join(build.Default.GOPATH, "include")

	MyKitBase = "trunkit"
	MyKitPath = filepath.Join(GoPath, MyKitBase)

	ValidateBase = "github.com/envoyproxy/protoc-gen-validate"
	ValidatePath = filepath.Join(GoPath, ValidateBase)

	MarketplaceBase = "marketplace"
	MarketplacePath = filepath.Join(GoPath, MarketplaceBase)
)
