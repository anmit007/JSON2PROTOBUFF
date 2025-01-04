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
