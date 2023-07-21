package redis

import (
	"context"
	"development/common"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var client *redis.Client

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	opts, err := redis.ParseURL(values.REDIS)
	if err != nil {
		log.Fatalln(err)
	}
	client = redis.NewClient(opts)
	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	d, _ := time.ParseDuration("600s")
	if err := client.Set(context.TODO(), "dev:example", 1234, d).Err(); err != nil {
		t.Error(err)
	}
}

func TestMock(t *testing.T) {
	ctx := context.TODO()
	x := client.TxPipeline()
	for i := 0; i < 10000; i++ {
		id, _ := gonanoid.New()
		x.Set(ctx, fmt.Sprintf(`session:%s`, id), funk.RandomString(10), 0)
	}
	if _, err := x.Exec(ctx); err != nil {
		t.Error(err)
	}
}

func TestScan(t *testing.T) {
	ctx := context.TODO()
	var keys []string
	var cursor uint64
	for {
		_keys, _cursor, err := client.Scan(ctx, cursor, "session:*", 100).Result()
		if err != nil {
			return
		}
		keys = append(keys, _keys...)
		if _cursor == 0 {
			break
		}
		cursor = _cursor
	}
	t.Log(keys)
	t.Log(len(keys))
}

var lkey = "l:123456"

func TestLPush(t *testing.T) {
	ctx := context.TODO()
	err := client.Del(ctx, lkey).Err()
	assert.NoError(t, err)
	err = client.LPush(ctx, lkey, time.Now().Format(time.RFC3339)).Err()
	assert.NoError(t, err)
	err = client.LPush(ctx, lkey, "asd").Err()
	assert.NoError(t, err)
	err = client.LPush(ctx, lkey, "acxc").Err()
	assert.NoError(t, err)
}

func TestRPop(t *testing.T) {
	ctx := context.TODO()
	begin, err := client.RPop(ctx, lkey).Time()
	assert.NoError(t, err)
	t.Log(begin)
	t.Log(time.Since(begin))
	t.Log(time.Since(begin) > time.Second*3)

	n, err := client.LLen(ctx, lkey).Result()
	assert.NoError(t, err)
	for n != 0 {
		v, err := client.RPop(ctx, lkey).Result()
		assert.NoError(t, err)
		t.Log(v)
		n--
	}
}

func TestSets(t *testing.T) {
	ctx := context.TODO()
	r, err := client.SAdd(ctx, "tx", "123456").Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestSave(t *testing.T) {
	ctx := context.TODO()
	r, err := client.Save(ctx).Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestLPush2(t *testing.T) {
	ctx := context.TODO()
	v := client.LPush(ctx, "test", time.Now().Format(time.RFC3339)).Val()
	t.Log(v)
}
