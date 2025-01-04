package json2protobuff

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ConvertSnakeCaseToUpperCamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = cases.Title(language.English).String(str)
	return strings.Replace(str, " ", "", -1)
}
func validateFieldName(fieldName string) bool {
	if len(fieldName) == 0 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	return re.MatchString(fieldName)
}

func (p *Parser) validName(raw *Message) error {

	for _, field := range raw.Entries {
		if !validateFieldName(field.Name) {
			return fmt.Errorf("invalid field name %s", field.Name)
		}
	}

	for mName, message := range raw.SubMessage {
		if !validateFieldName(mName) {
			return fmt.Errorf("invalid field name %s", mName)
		}
		for _, field := range message.Entries {
			if !validateFieldName(field.Name) {
				return fmt.Errorf("invalid field name %s.%s", mName, field.Name)
			}
		}
	}
	return nil
}
func GoType2ProtoType(to reflect.Kind) string {
	switch to {
	case reflect.Int8, reflect.Int16, reflect.Int32:
		return "sint32"
	case reflect.Int, reflect.Int64:
		return "sint64"
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "fixed32"
	case reflect.Uint, reflect.Uint64:
		return "fixed64"
	case reflect.Bool:
		return "bool"
	case reflect.String:
		return "string"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	default:
		panic("invalid type: " + to.String())
	}
}
