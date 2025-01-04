package json2protobuff

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"gitter.top/common/goref"
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

func WithTiledMessage(isTiled bool) func(parser *Parser) {
	return func(parser *Parser) {
		parser.TiledDefinition = isTiled
	}
}

func WithMergeMessage(isMerge bool) func(parser *Parser) {
	return func(parser *Parser) {
		parser.MergeMessage = isMerge
	}
}

func WithShowFileHeader(isShow bool) func(parser *Parser) {
	return func(parser *Parser) {
		parser.FileHeaderDefinition = isShow
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

func (p *Parser) rangeItem(mapIter *reflect.MapIter) *Message {
	if mapIter == nil {
		return nil
	}
	var msg = &Message{SubMessage: make(map[string]*Message)}
	for mapIter.Next() {
		var fd = new(Entry)
		fd.Name = mapIter.Key().String()
		value := mapIter.Value()
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if value.Kind() == reflect.Interface {
			value = reflect.ValueOf(value.Interface())
		}
		switch value.Kind() {
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
			fd.Type = GoType2ProtoType(value.Kind())
		case reflect.Map:
			fd.Type = ConvertSnakeCaseToUpperCamelCase(fd.Name)
			m := p.rangeItem(value.MapRange())
			msg.SubMessage[fd.Type] = m
		case reflect.Slice, reflect.Array:
			fd.Repeated = true
			size := value.Len()
			if size == 0 {
				fd.Type = "google.protobuf.Any"
			} else {
				vi := value.Index(0).Interface()
				if goref.IsBasicType(vi) {
					fd.Type = GoType2ProtoType(reflect.ValueOf(vi).Type().Kind())
				} else {
					fd.Type = ConvertSnakeCaseToUpperCamelCase(fd.Name)
					m := p.rangeItem(reflect.ValueOf(vi).MapRange())
					msg.SubMessage[fd.Type] = m
				}
			}
		case reflect.Invalid:
			fd.Type = "google.protobuf.Any"
		default:
			fmt.Printf("Unhandled type - key: %s, type: %s, value: %v\n",
				fd.Name, value.Kind().String(), mapIter.Value())
		}

		msg.Entries = append(msg.Entries, fd)
	}
	return msg
}
