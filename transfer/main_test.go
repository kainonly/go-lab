package transfer

import (
	"context"
	"development/common"
	natscommon "development/nats/common"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/transfer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var client *mongo.Client
var db *mongo.Database
var x *transfer.Transfer

func TestMain(m *testing.M) {
	var err error
	path := "../config/config.yml"
	if values, err = common.LoadValues(path); err != nil {
		log.Fatalln(err)
	}
	if client, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(values.MONGO),
	); err != nil {
		log.Fatalln(err)
	}

	option := options.Database().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	db = client.Database("xapi", option)

	var nc *nats.Conn
	if nc, err = natscommon.Create(path); err != nil {
		log.Fatalln(err)
	}
	var js nats.JetStreamContext
	if js, err = nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
		log.Fatalln(err)
	}
	if x, err = transfer.New("xapi", db, js); err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	if err := x.Set(context.TODO(), transfer.Option{
		Key:         "test",
		Description: "测试",
		TTL:         3600,
	}); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	result, err := x.Get("test")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestPublish(t *testing.T) {
	if err := x.Publish(context.TODO(), "test", transfer.Payload{
		Metadata: map[string]interface{}{
			"uuid": "0ff5483a-7ddc-44e0-b723-c3417988663f",
		},
		Data: map[string]interface{}{
			"msg": "hi",
		},
		Timestamp: time.Now(),
	}); err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	if err := x.Remove("test"); err != nil {
		t.Error(err)
	}
}
