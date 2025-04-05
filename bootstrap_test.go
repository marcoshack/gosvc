package gosvc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcoshack/gosvc"
	gsvtesting "github.com/marcoshack/gosvc/internal/testing"
)

func TestBootstrap_LoadsConfiguration(t *testing.T) {
	bs, err := gosvc.NewBootstrap[gsvtesting.TestConfigType](context.Background(), gosvc.BootstrapInput{
		Name:           "TestService",
		ConfigFileName: "internal/testing/data/config.json",
	})
	require.NoError(t, err)
	require.NotNil(t, bs)
	require.NotNil(t, bs.Config)
	assert.Equal(t, "value1", bs.Config.Attr1)
	assert.Equal(t, 123, bs.Config.Attr2)
	assert.Equal(t, "subValue1", bs.Config.Attr5.SubAttr1)
	assert.Contains(t, bs.Config.Attr6, "listValue1")
	assert.Contains(t, bs.Config.Attr6, "listValue2")
}
