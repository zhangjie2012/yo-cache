package cache

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhangjie2012/yo-kit/datastruct"
)

func TestMain(m *testing.M) {
	c := datastruct.RedisConf{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}
	if err := Init(context.Background(), "yo-cache", &c); err != nil {
		log.Fatal(err)
	}
	defer Close()

	m.Run()
}

func TestTTL(t *testing.T) {
	var (
		ctx context.Context = context.Background()
		key string          = "ttl.test"
		err error
		d   time.Duration
	)

	d, err = TTL(ctx, key)
	require.Nil(t, err)
	assert.Equal(t, -2*time.Nanosecond, d)
}

func TestPTTL(t *testing.T) {
	var (
		ctx context.Context = context.Background()
		key string          = "ttl.test"
		err error
		d   time.Duration
	)

	d, err = PTTL(ctx, key)
	require.Nil(t, err)
	assert.Equal(t, -2*time.Nanosecond, d)
}

func TestSetGetObject(t *testing.T) {
	type User struct {
		Name    string
		Email   string
		Address string
	}

	var (
		ctx   context.Context = context.Background()
		key   string          = "get-set.object.user"
		user                  = User{}
		user2                 = User{}
		err   error
	)

	err = GetObject(ctx, key, &user)
	assert.Equal(t, ErrNotExist, err)

	user.Name = "JerrZhang"
	user.Email = "me@zhangjiee.com"
	user.Address = "china"

	err = SetObject(ctx, key, &user, 30*time.Second)
	require.Nil(t, err)

	err = GetObject(ctx, key, &user2)
	require.Nil(t, err)
	assert.EqualValues(t, user2.Name, user.Name)

	err = Del(ctx, key)
	require.Nil(t, err)
}

func TestGetSetString(t *testing.T) {
	var (
		ctx   context.Context = context.Background()
		key   string          = "get-set.string"
		user                  = "hello, world"
		user2                 = ""
		err   error
	)

	err = SetString(ctx, key, user, 1*time.Minute)
	require.Nil(t, err)

	user2, err = GetString(ctx, key)
	require.Nil(t, err)

	assert.Equal(t, user, user2)

	Del(ctx, key)
}

func TestGetSetInt(t *testing.T) {
	var (
		ctx    context.Context = context.Background()
		key    string          = "get-set.int"
		value1 int             = 1024
		value2 int             = 0
		err    error
	)

	err = SetInt(ctx, key, value1, 1*time.Minute)
	require.Nil(t, err)

	value2, err = GetInt(ctx, key)
	require.Nil(t, err)
	assert.Equal(t, value1, value2)

	Del(ctx, key)
}

func TestGetSetInt64(t *testing.T) {
	var (
		ctx    context.Context = context.Background()
		key    string          = "get-set.int"
		value1 int64           = 102410241024
		value2 int64           = 0
		err    error
	)

	err = SetInt64(ctx, key, value1, 1*time.Minute)
	require.Nil(t, err)

	value2, err = GetInt64(ctx, key)
	require.Nil(t, err)
	assert.Equal(t, value1, value2)

	Del(ctx, key)
}

func TestGetSetFloat64(t *testing.T) {
	var (
		ctx    context.Context = context.Background()
		key    string          = "get-set.float64"
		value1 float64         = 10241024.1024
		value2 float64         = 0
		err    error
	)

	err = SetFloat64(ctx, key, value1, 1*time.Minute)
	require.Nil(t, err)

	value2, err = GetFloat64(ctx, key)
	require.Nil(t, err)
	assert.Equal(t, value1, value2)

	Del(ctx, key)
}

func TestGetSetBool(t *testing.T) {
	var (
		ctx    context.Context = context.Background()
		key    string          = "get-set.bool"
		value1 bool            = true
		value2 bool            = false
		err    error
	)

	err = SetBool(ctx, key, value1, 1*time.Minute)
	require.Nil(t, err)

	value2, err = GetBool(ctx, key)
	require.Nil(t, err)
	assert.Equal(t, value1, value2)

	Del(ctx, key)
}
