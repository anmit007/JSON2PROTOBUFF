package json2protobuff

import (
	"regexp"
)

type Entry struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Repeated bool   `json:"repeated,omitempty"`
}

type Message struct {
	SubMessage map[string]*Message `json:"sub_message,omitempty"`
	Entries    []*Entry            `json:"entries,omitempty"`
}

func validateFieldName(fieldName string) bool {

	if len(fieldName) == 0 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	return re.MatchString(fieldName)
}
