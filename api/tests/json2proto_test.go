package tests

import (
	"strings"
	"testing"

	json2protobuff "github.com/anmit007/JSON2PROTOBUFF"
	"github.com/stretchr/testify/assert"
)

func TestNewParser_BasicInit(t *testing.T) {
	var test = `{"name":"test","age":25}`
	reader := strings.NewReader(test)

	parser, err := json2protobuff.NewParser(reader,
		json2protobuff.WithShowFileHeader(true),
		json2protobuff.WithMergeMessage(true),
		json2protobuff.WithTiledMessage(false))
	assert.NoError(t, err)
	assert.NotNil(t, parser)
	assert.True(t, parser.FileHeaderDefinition)
	assert.True(t, parser.MergeMessage)
	assert.False(t, parser.TiledDefinition)
}

func TestTemplateRendering(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tiled    bool
		expected string
	}{
		{
			name:  "simple object tiled",
			input: `{"name": "test", "age": 25}`,
			tiled: true,
			expected: `message AutoGenerated {
    string name = 1;
    sint32 age = 2;
}`,
		},
		{
			name:  "nested object embedded",
			input: `{"user": {"name": "test", "age": 25}}`,
			tiled: false,
			expected: `message AutoGenerated {
    message User {
        string name = 1;
        sint32 age = 2;
    }
    User user = 1;
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			parser, err := json2protobuff.NewParser(reader,
				json2protobuff.WithTiledMessage(tt.tiled),
				json2protobuff.WithShowFileHeader(false))
			assert.NoError(t, err)

			err = parser.Parse()
			assert.NoError(t, err)

			output := parser.Output()

			cleanOutput := strings.Join(strings.Fields(output), " ")
			cleanExpected := strings.Join(strings.Fields(tt.expected), " ")
			assert.Equal(t, cleanExpected, cleanOutput)
		})
	}
}
