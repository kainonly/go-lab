package main

import (
	"development/common"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	values, err := common.LoadValues("../../config.yml")
	if err != nil {
		panic(err)
	}
	if err = UseNats(values); err != nil {
		panic(err)
	}
	if js, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
	); err != nil {
		return
	}
	os.Exit(m.Run())
}

func TestPublic(t *testing.T) {
	_, err := js.Publish("test", []byte("x2"))
	assert.NoError(t, err)
}
