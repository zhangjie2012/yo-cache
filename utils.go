package cache

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

// ComposeKey(2) 合成 key, 模块使用 `:` 分隔, 子模块使用 `.`
func ComposeKey(rawkey string) string {
	return fmt.Sprintf("%s:%s:%s", version_, appName_, rawkey)
}
func ComposeKey2(module, rawkey string) string {
	return fmt.Sprintf("%s:%s:%s:%s", version_, appName_, module, rawkey)
}

// ErrorWrapper 错误码劫持, 全局扩展
func ErrorWrapper(err error) error {
	if err == redis.Nil {
		return ErrNotExist
	}
	return err
}
