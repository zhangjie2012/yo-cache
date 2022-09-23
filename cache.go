package cache

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zhangjie2012/yo-kit/datastruct"
)

var (
	ErrNotExist       = errors.New("ERROR: key not exist")
	ErrInitRepeat     = errors.New("ERROR: init repeat")
	ErrAppNameInvalid = errors.New("ERROR: invalid appname")
)

var (
	prefix_  string        = "yo_cache"
	appName_ string        = ""
	version_ string        = "v1" // 用于 package 不兼容更新
	client_  *redis.Client = nil
)

// Init 初始化 cache, appName 用于 key 隔离
func Init(ctx context.Context, appName string, c *datastruct.RedisConf) error {
	appName = strings.TrimSpace(appName)
	if appName == "" || strings.Index(appName, ":") != -1 {
		return ErrAppNameInvalid
	}

	if client_ != nil {
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	appName_ = appName
	client_ = client

	return nil
}

func Close() (err error) {
	if client_ != nil {
		if err = client_.Close(); err != nil {
			return
		}
		client_ = nil
	}
	return
}

// C expose redis client for native redis library visit
func C() *redis.Client {
	return client_
}

// TTL seconds resolution
// - The command returns -1 if the key exists but has no associated expire.
// - The command returns -2 if the key does not exist.
func TTL(ctx context.Context, key string) (time.Duration, error) {
	realKey := ComposeKey(key)
	d, err := client_.TTL(ctx, realKey).Result()
	if err != nil {
		return 0, ErrorWrapper(err)
	}
	return d, nil
}

// PTTL milliseconds resolution
// - The command returns -1 if the key exists but has no associated expire.
// - The command returns -2 if the key does not exist.
func PTTL(ctx context.Context, key string) (time.Duration, error) {
	realKey := ComposeKey(key)
	d, err := client_.PTTL(ctx, realKey).Result()
	if err != nil {
		return 0, ErrorWrapper(err)
	}
	return d, nil
}

func Del(ctx context.Context, key string) error {
	realKey := ComposeKey(key)
	err := client_.Del(ctx, realKey).Err()
	if err != nil {
		return ErrorWrapper(err)
	}
	return nil
}

// Set/Get object
func SetObject(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	realKey := ComposeKey(key)

	bs, err := msgpack.Marshal(value)
	if err != nil {
		return ErrorWrapper(err)
	}

	if err := client_.Set(ctx, realKey, bs, expire).Err(); err != nil {
		return ErrorWrapper(err)
	}

	return nil
}
func GetObject(ctx context.Context, key string, value interface{}) error {
	realKey := ComposeKey(key)

	bs, err := client_.Get(ctx, realKey).Bytes()
	if err != nil {
		return ErrorWrapper(err)
	}

	if err := msgpack.Unmarshal(bs, value); err != nil {
		return ErrorWrapper(err)
	}

	return nil
}

func SetString(ctx context.Context, key string, value string, expire time.Duration) error {
	realKey := ComposeKey(key)

	err := client_.Set(ctx, realKey, value, expire).Err()
	if err != nil {
		return ErrorWrapper(err)
	}

	return nil
}

func GetString(ctx context.Context, key string) (string, error) {
	realKey := ComposeKey(key)

	value, err := client_.Get(ctx, realKey).Result()
	if err != nil {
		return "", ErrorWrapper(err)
	}

	return value, err
}

func SetInt(ctx context.Context, key string, value int, expire time.Duration) error {
	return SetString(ctx, key, strconv.Itoa(value), expire)
}

func GetInt(ctx context.Context, key string) (int, error) {
	realKey := ComposeKey(key)

	value, err := client_.Get(ctx, realKey).Int()
	if err != nil {
		return 0, ErrorWrapper(err)
	}

	return value, err
}

func SetInt64(ctx context.Context, key string, value int64, expire time.Duration) error {
	return SetString(ctx, key, strconv.FormatInt(value, 10), expire)
}

func GetInt64(ctx context.Context, key string) (int64, error) {
	realKey := ComposeKey(key)

	value, err := client_.Get(ctx, realKey).Int64()
	if err != nil {
		return 0, ErrorWrapper(err)
	}

	return value, err
}

func SetFloat64(ctx context.Context, key string, value float64, expire time.Duration) error {
	return SetString(ctx, key, strconv.FormatFloat(value, 'f', -1, 64), expire)
}

func GetFloat64(ctx context.Context, key string) (float64, error) {
	realKey := ComposeKey(key)

	value, err := client_.Get(ctx, realKey).Float64()
	if err != nil {
		return 0, ErrorWrapper(err)
	}

	return value, err
}

func SetBool(ctx context.Context, key string, b bool, expire time.Duration) error {
	if b {
		return SetInt(ctx, key, 1, expire)
	}
	return SetInt(ctx, key, 0, expire)
}

func GetBool(ctx context.Context, key string) (bool, error) {
	value, err := GetInt(ctx, key)
	if err != nil {
		return false, err
	}
	if value == 1 {
		return true, err
	}
	return false, err
}

// // -----------------------------------------------------------------------------
// // Set wrapper
// // SS for Set String
// // -----------------------------------------------------------------------------

// // SSMembers get all members slice
// func SSMembers(key string) ([]string, error) {
// 	aKey := composeKey2(setModule, key)
// 	values, err := redisClient.SMembers(aKey).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			return nil, NotExist
// 		}
// 		return nil, err
// 	}
// 	return values, nil
// }

// // SSAdd add members to Set
// func SSAdd(key string, members ...string) error {
// 	aKey := composeKey2(setModule, key)
// 	t := []interface{}{}
// 	for _, v := range members {
// 		t = append(t, v)
// 	}
// 	_, err := redisClient.SAdd(aKey, t...).Result()
// 	return err
// }

// // SSRem remove members from Set
// func SSRem(key string, members ...string) error {
// 	aKey := composeKey2(setModule, key)
// 	t := []interface{}{}
// 	for _, v := range members {
// 		t = append(t, v)
// 	}
// 	_, err := redisClient.SRem(aKey, t...).Result()
// 	return err
// }

// // SSCount get member count
// func SSCount(key string) int64 {
// 	aKey := composeKey2(setModule, key)
// 	count, err := redisClient.SCard(aKey).Result()
// 	if err != nil {
// 		return 0
// 	}
// 	return count
// }

// // SSIsMember check set if include member
// func SSIsMember(key string, member string) bool {
// 	aKey := composeKey2(setModule, key)
// 	ok, err := redisClient.SIsMember(aKey, member).Result()
// 	if err != nil {
// 		return false
// 	}
// 	return ok
// }

// // SSRandomN random get N members
// func SSRandomN(key string, count int64) []string {
// 	aKey := composeKey2(setModule, key)
// 	values, err := redisClient.SRandMemberN(aKey, count).Result()
// 	if err != nil {
// 		return []string{}
// 	}
// 	return values
// }

// func SSDelete(key string) {
// 	aKey := composeKey2(setModule, key)
// 	redisClient.Del(aKey)
// }

// // SS_TTL seconds resolution
// // - The command returns -1 if the key exists but has no associated expire.
// // - The command returns -2 if the key does not exist.
// func SS_TTL(key string) time.Duration {
// 	aKey := composeKey2(setModule, key)
// 	d, err := redisClient.TTL(aKey).Result()
// 	if err != nil {
// 		return 0
// 	}
// 	return d
// }

// // SS_TTL milliseconds resolution
// // - The command returns -1 if the key exists but has no associated expire.
// // - The command returns -2 if the key does not exist.
// func SS_PTTL(key string) time.Duration {
// 	aKey := composeKey2(setModule, key)
// 	d, err := redisClient.PTTL(aKey).Result()
// 	if err != nil {
// 		return 0
// 	}
// 	return d
// }

// func SSExpire(key string, d time.Duration) error {
// 	aKey := composeKey2(setModule, key)
// 	return redisClient.Expire(aKey, d).Err()
// }
