package init

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mykit/internal/constant"
	"mykit/internal/emitter/common"
	"mykit/internal/metadata"
	osutil "mykit/internal/util/os"
)

var _initFilePaths = []struct {
	Src             string
	Dest            string
	SkipForMonorepo bool
}{
	{
		Src:             "goservice/api/proto.tmpl",
		Dest:            "api/xtype.proto",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/api/config.proto.tmpl",
		Dest:            "api/xtype_config.proto",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/api/code.proto.tmpl",
		Dest:            "api/xtype_code.proto",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/configs/config.yaml.tmpl",
		Dest:            "configs/config.yaml",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/schema/sample.go.tmpl",
		Dest:            "schema/sample.go",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/schema/base.go.tmpl",
		Dest:            "schema/base.go",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/mykit.yaml.tmpl",
		Dest:            "mykit.yaml",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/Makefile.tmpl",
		Dest:            "Makefile",
		SkipForMonorepo: false,
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
		Src:             "goservice/ci.yaml.tmpl",
		Dest:            "ci.yaml",
		SkipForMonorepo: false,
	},
	{
		Src:             "goservice/gitlab-ci.yml.tmpl",
		Dest:            ".gitlab-ci.yml",
		SkipForMonorepo: true,
	},
	{
		Src:             "goservice/go.mod.tmpl",
		Dest:            "go.mod",
		SkipForMonorepo: true,
	},
}

func Init(pkg, name string, monorepo bool, typ string) {
	name = strings.ToLower(name)

	fmt.Println("New Project:")
	fmt.Println("- Name:", name)
	fmt.Println("- Directory:", metadata.Dir)
	fmt.Println("- Package:", pkg)
	fmt.Println("- Monorepo:", monorepo)
	fmt.Println()

	fmt.Println("Create files:")
	for _, p := range _initFilePaths {
		dest := strings.ReplaceAll(p.Dest, "xtype", name)
		destPath := filepath.Join(metadata.Dir, dest)

		if err := osutil.CreateDirIfNotExists(filepath.Dir(destPath), 0755); err != nil {
			os.Exit(1)
		}

		if monorepo && p.SkipForMonorepo {
			continue
		}

		fmt.Println("- Create", destPath)

		protoGoPackage := name
		if strings.HasSuffix(metadata.Dir, pkg) {
			protoGoPackage = pkg
		} else if strings.HasPrefix(metadata.Dir, constant.GoPath) {
			protoGoPackage = strings.TrimPrefix(metadata.Dir, constant.GoPath+string(filepath.Separator))
		}

		common.Render(p.Src, destPath, map[string]interface{}{
			"MyKitBase":      constant.MyKitBase,
			"ProjectName":    name,
			"ServiceName":    name,
			"Package":        pkg,
			"Version":        metadata.MyKitVersion,
			"ProtoGoPackage": protoGoPackage,

			// below are for mykit.yaml only
			"Monorepo": monorepo,
			"Ent":      false,
			"Client":   true,
			"Type":     typ,
			"GenerateGo": []string{
				fmt.Sprintf("%s.proto", name),
				fmt.Sprintf("%s_config.proto", name),
				fmt.Sprintf("%s_code.proto", name),
			},
			"GenerateJs": []string{
				fmt.Sprintf("%s.proto", name),
				fmt.Sprintf("%s_code.proto", name),
			},
		})
	}
	fmt.Println()

	cmd := exec.Command("mykit", "generate", "go", "-u")
	cmd.Dir = metadata.Dir
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("generate go failed", err)
		os.Exit(1)
	}

	err = exec.Command("gofmt", "-w", metadata.Dir).Run()
	if err != nil {
		fmt.Println("gofmt failed", metadata.Dir, err)
		os.Exit(1)
	}

	fmt.Println("Init", name, "done!")
}
