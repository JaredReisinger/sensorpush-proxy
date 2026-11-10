package cmd

import (
	"testing"

	"github.com/jaredreisinger/asp"
	"github.com/stretchr/testify/assert"
)

// This is a bogus test just meant to place-hold for future actual cobra command
// testing.
func TestProxyConfig(t *testing.T) {
	s, err := asp.SerializeFlags(proxyConfig{}, false)
	assert.Nil(t, err)
	assert.Equal(t, `--port "" --sensors "" --update-period ""`, s)
}
