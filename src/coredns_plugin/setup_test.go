package squawker

import (
	"testing"

	"github.com/coredns/caddy"
	"github.com/stretchr/testify/assert"
)

func TestSetupPlugin(t *testing.T) {
	config, err := parseArgs(caddy.NewTestController("dns", "squawker"))
	assert.Nil(t, config, "config not returned on error")
	assert.Error(t, err, "error on no args")

	config, err = parseArgs(caddy.NewTestController("dns", "squawker 192.168.1.1:8080"))
	assert.NotNil(t, config, "config returned on provided address")
	assert.Equal(t, "192.168.1.1:8080", config.address, "match provided address")
	assert.NoError(t, err, "no error on provided address")
}
