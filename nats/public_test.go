package nats

import (
	"github.com/nats-io/nats.go"
	"sync"
	"testing"
)

func TestPublishN(t *testing.T) {
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
