package template

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

func Lowercase(name string) string {
	return strings.ToLower(name)
}

func Camelcase(name string) string {
	return strcase.ToCamel(name)
}

func Uppercase(name string) string {
	return strings.ToUpper(name)
}

func LowerCamel(name string) string {
	return strcase.ToLowerCamel(name)
}

func Last(x int, a interface{}) bool {
	return x == reflect.ValueOf(a).Len()-1
}

func CamelToSnake(s string) string {
	return strcase.ToSnake(s)
}
