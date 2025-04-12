package bootstrap_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcoshack/gosvc/bootstrap"
	gsvtesting "github.com/marcoshack/gosvc/internal/testing"
)

var (
	testArgs = []string{"-c", "../internal/testing/data/config.json"}
)

func TestBootstrap_LoadsConfiguration(t *testing.T) {
	bs, err := bootstrap.New[gsvtesting.TestConfigType](context.Background(), bootstrap.Input{
		ServiceName: "TestService",
		AWSRegion:   "us-east-1",
		Args:        testArgs,
	})
	require.NoError(t, err)
	require.NotNil(t, bs)

	assert.NotNil(t, bs.Logger)
	assert.NotNil(t, bs.Ctx)

	require.NotNil(t, bs.Config)
	assert.Equal(t, "value1", bs.Config.Attr1)
	assert.Equal(t, 123, bs.Config.Attr2)
	assert.Equal(t, "subValue1", bs.Config.Attr5.SubAttr1)
	assert.Contains(t, bs.Config.Attr6, "listValue1")
	assert.Contains(t, bs.Config.Attr6, "listValue2")
}
