package redis

import (
	"context"
	"development/common"
	"fmt"
	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
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