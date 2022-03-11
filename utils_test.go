package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeKey(t *testing.T) {
	appName_ = "app"

	key1 := ComposeKey("hello")
	assert.Equal(t, "v1:app:hello", key1)

	key2 := ComposeKey2("module", "hello")
	assert.Equal(t, "v1:app:module:hello", key2)
}
