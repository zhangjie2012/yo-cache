package cache

import (
	"strings"

	"github.com/go-redis/redis/v8"
)

// ComposeKey 合成 key, 模块使用 `:` 分隔, 子模块使用 `.`
func ComposeKey(rawkey string) string {
	return ComposeKeys(prefix_, version_, appName_, rawkey)
}

// ErrorWrapper 错误码劫持, 全局扩展
func ErrorWrapper(err error) error {
	if err == redis.Nil {
		return ErrNotExist
	}
	return err
}

// ComposeKeys
func ComposeKeys(keys ...string) string {
	return strings.Join(keys, ":")
}
