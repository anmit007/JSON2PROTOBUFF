package json2protobuff

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ConvertSnakeCaseToUpperCamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = cases.Title(language.English).String(str)
	return strings.Replace(str, " ", "", -1)
}
