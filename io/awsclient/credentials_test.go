package awsclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCredentials(t *testing.T) {
	assert.True(t, credentials.HasKeys())
	assert.False(t, credentials.Expired())
}
