package template

import (
	"bytes"
	"fmt"
	"text/template"
)

type Data interface{}

func Generate(buffer *bytes.Buffer, data Data, content string) error {
	funcMap := template.FuncMap{
		"uppercase":  Uppercase,
		"lowercase":  Lowercase,
		"camelcase":  Camelcase,
		"lowerCamel": LowerCamel,
		"last":       Last,
	}
	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(content)
	if err != nil {
		return fmt.Errorf("error when parse template %v", err)
	}
	err = tmpl.Execute(buffer, data)
	if err != nil {
		return fmt.Errorf("error when exec template %v", err)
	}

	return nil
}
