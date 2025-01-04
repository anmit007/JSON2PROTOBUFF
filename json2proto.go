package json2protobuff

import (
	"encoding/json"
	"fmt"
	"io"
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

type Parser struct {
	FileHeaderDefinition bool
	TiledDefinition      bool
	MergeMessage         bool
	_data                map[string]interface{}
	_rawMessage          *Message
	_tiledMessage        *Message
}
type Option func(Parser *Parser)

func WithTiledMessage(isTiled bool) Option {
	return func(parser *Parser) {
		parser.TiledDefinition = isTiled
	}
}

func NewParser(reader io.Reader, options ...Option) (*Parser, error) {
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if !json.Valid(buffer) {
		return nil, fmt.Errorf("ivalid json format")
	}
	parser := &Parser{
		_data:         make(map[string]interface{}),
		_rawMessage:   &Message{SubMessage: make(map[string]*Message)},
		_tiledMessage: &Message{SubMessage: make(map[string]*Message)},
	}
	for _, option := range options {
		option(parser)
	}
	if err := json.Unmarshal(buffer, &parser._data); err != nil {
		return nil, err
	}
	return parser, nil
}

func validateFieldName(fieldName string) bool {

	if len(fieldName) == 0 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	return re.MatchString(fieldName)
}
