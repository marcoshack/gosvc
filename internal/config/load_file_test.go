package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcoshack/gosvc/internal/config"
	gsvtesting "github.com/marcoshack/gosvc/internal/testing"
)

func TestConfig_LoadFromfile(t *testing.T) {
	config, err := config.LoadFromFile[gsvtesting.TestConfigType](&config.LoadFromFileInput{FileName: "../testing/data/config.json"})
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "value1", config.Attr1)
	assert.Equal(t, 123, config.Attr2)
	assert.Equal(t, "subValue1", config.Attr5.SubAttr1)
	assert.Contains(t, config.Attr6, "listValue1")
	assert.Contains(t, config.Attr6, "listValue2")
}
