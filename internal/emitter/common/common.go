package common

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"mykit/internal/metadata"
	"mykit/internal/template"
	"mykit/internal/templates"
	osutil "mykit/internal/util/os"
)

func Render(templatePath string, outPath string, data template.Data, options ...RenderOption) {
	renderOptions := &renderOptions{}
	for _, o := range options {
		o.Apply(renderOptions)
	}

	if !strings.HasPrefix(outPath, metadata.Dir) {
		outPath = filepath.Join(metadata.Dir, outPath)
	}
	if !renderOptions.overwrite {
		if _, err := os.Stat(outPath); err == nil {
			return
		}
	}

	if err := osutil.CreateDirIfNotExists(filepath.Dir(outPath), 0755); err != nil {
		os.Exit(1)
	}

	templateContent, err := templates.Load(templatePath)
	if err != nil {
		fmt.Println(fmt.Sprintf("load template failed %s", templatePath), err)
		os.Exit(1)
	}
	buffer := new(bytes.Buffer)
	err = template.Generate(buffer, data, string(templateContent))
	if err != nil {
		fmt.Println(fmt.Sprintf("template generate failed %s", templatePath), err)
		os.Exit(1)
	}

	p := buffer.Bytes()
	if strings.HasSuffix(outPath, ".go") {
		f, err := format.Source(buffer.Bytes())
		if err != nil {
			fmt.Println(fmt.Sprintf("format failed %s", templatePath), err)
			os.Exit(1)
		}
		p = f
	}

	err = ioutil.WriteFile(outPath, p, 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("write failed %s", outPath), err)
		os.Exit(1)
	}
}

type RenderOption interface {
	Apply(o *renderOptions)
}

type renderOptions struct {
	overwrite bool
}

type OptionFunc func(*renderOptions)

func (f OptionFunc) Apply(o *renderOptions) {
	f(o)
}

func Overwrite() RenderOption {
	return OptionFunc(func(o *renderOptions) {
		o.overwrite = true
	})
}
