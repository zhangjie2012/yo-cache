package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeKey(t *testing.T) {
	appName_ = "app"

	key1 := ComposeKey("hello")
	assert.Equal(t, "yo_cache:v1:app:hello", key1)
}
