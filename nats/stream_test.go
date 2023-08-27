package nats

import (
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStreamNames(t *testing.T) {
	for name := range js.StreamNames() {
		t.Log(name)
	}
}

func TestStreamsInfo(t *testing.T) {
	for v := range js.Streams() {
		data, _ := sonic.Marshal(v)
		t.Log(string(data))
	}
}

func TestDeleteStreams(t *testing.T) {
	for name := range js.StreamNames() {
		err := js.DeleteStream(name)
		assert.NoError(t, err)
	}
}

func TestAddStream(t *testing.T) {
	info, err := js.AddStream(&nats.StreamConfig{
		Name:      "development",
		Subjects:  []string{"development"},
		Retention: nats.WorkQueuePolicy,
	})
	assert.NoError(t, err)
	t.Log(info)
}

func TestPublish(t *testing.T) {
	_, err := js.Publish("development.message", []byte("hello"))
	assert.NoError(t, err)
}

func TestStreamInfo(t *testing.T) {
	v, err := js.StreamInfo("development")
	assert.NoError(t, err)
	data, _ := sonic.Marshal(v)
	t.Log(string(data))
}

func TestUpdateStream(t *testing.T) {
	info, err := js.UpdateStream(&nats.StreamConfig{
		Name:        "development",
		Subjects:    []string{"development.>"},
		Description: "测试",
		Retention:   nats.WorkQueuePolicy,
	})
	assert.NoError(t, err)
	t.Log(info)
}

func TestDeleteStream(t *testing.T) {
	err := js.DeleteStream("development")
	assert.NoError(t, err)
}
