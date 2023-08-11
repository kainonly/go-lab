package nats

import (
	"development/common"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sync"
	"testing"
)

var nc *nats.Conn
var js nats.JetStreamContext

func TestMain(m *testing.M) {
	values, err := common.LoadValues("./config.yml")
	if err != nil {
		panic(err)
	}
	if nc, err = common.UseNats(values); err != nil {
		panic(err)
	}
	if js, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
	); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestPublishMessage(t *testing.T) {
	err := nc.Publish("development.message", []byte("abc"))
	assert.NoError(t, err)
}

func TestPublishJs(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		js.Publish("development", []byte("abc"), nats.MsgId("abc"))
		wg.Done()
	}()
	go func() {
		js.Publish("development", []byte("abC"), nats.MsgId("abx"))
		wg.Done()
	}()
	go func() {
		js.Publish("development", []byte("Abc"), nats.MsgId("aba"))
		wg.Done()
	}()
	wg.Wait()
}
