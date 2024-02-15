package redis

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/goutil/strutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golab/common"
	"os"
	"testing"
	"time"
)

var values *common.Values
var db *redis.Client

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		panic(err)
	}
	opts, err := redis.ParseURL(values.REDIS)
	if err != nil {
		panic(err)
	}
	db = redis.NewClient(opts)
	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	d, _ := time.ParseDuration("600s")
	if err := db.Set(context.TODO(), "dev:example", 1234, d).Err(); err != nil {
		t.Error(err)
	}
}

func TestMock(t *testing.T) {
	ctx := context.TODO()
	x := db.TxPipeline()
	for i := 0; i < 10000; i++ {
		id := strutil.MicroTimeHexID()
		x.Set(ctx, fmt.Sprintf(`session:%s`, id), strutil.RandomCharsV3(6), 0)
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
		_keys, _cursor, err := db.Scan(ctx, cursor, "session:*", 100).Result()
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
	err := db.Del(ctx, lkey).Err()
	assert.NoError(t, err)
	err = db.LPush(ctx, lkey, time.Now().Format(time.RFC3339)).Err()
	assert.NoError(t, err)
	err = db.LPush(ctx, lkey, "asd").Err()
	assert.NoError(t, err)
	err = db.LPush(ctx, lkey, "acxc").Err()
	assert.NoError(t, err)
}

func TestRPop(t *testing.T) {
	ctx := context.TODO()
	begin, err := db.RPop(ctx, lkey).Time()
	assert.NoError(t, err)
	t.Log(begin)
	t.Log(time.Since(begin))
	t.Log(time.Since(begin) > time.Second*3)

	n, err := db.LLen(ctx, lkey).Result()
	assert.NoError(t, err)
	for n != 0 {
		v, err := db.RPop(ctx, lkey).Result()
		assert.NoError(t, err)
		t.Log(v)
		n--
	}
}

func TestSets(t *testing.T) {
	ctx := context.TODO()
	r, err := db.SAdd(ctx, "tx", "123456").Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestSave(t *testing.T) {
	ctx := context.TODO()
	r, err := db.Save(ctx).Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestLPush2(t *testing.T) {
	ctx := context.TODO()
	v := db.LPush(ctx, "test", time.Now().Format(time.RFC3339)).Val()
	t.Log(v)
}
