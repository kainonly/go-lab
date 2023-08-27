package nats

import (
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsumerNames(t *testing.T) {
	for name := range js.ConsumerNames("development") {
		t.Log(name)
	}
}

func TestConsumersInfo(t *testing.T) {
	for v := range js.Consumers("development") {
		data, _ := sonic.Marshal(v)
		t.Log(string(data))
	}
}

func TestAddConsumer(t *testing.T) {
	info, err := js.AddConsumer("development", &nats.ConsumerConfig{
		Durable: "DEV",
	})
	assert.NoError(t, err)
	t.Log(info)
}

func TestConsumerInfo(t *testing.T) {
	info, err := js.ConsumerInfo("development", "DEV")
	assert.NoError(t, err)
	data, _ := sonic.Marshal(info)
	t.Log(string(data))
}
