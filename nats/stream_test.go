package nats

import (
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"testing"
)

func TestStreamNames(t *testing.T) {
	for name := range js.StreamNames() {
		t.Log(name)
	}
}

func TestStreamsInfo(t *testing.T) {
	//for v := range js.StreamsInfo() {
	//	data, _ := jsoniter.Marshal(v)
	//	t.Log(string(data))
	//}
	for v := range js.Streams() {
		data, _ := sonic.Marshal(v)
		t.Log(string(data))
	}
}

// TODO:删除所有 Stream
func TestDeleteStreams(t *testing.T) {
	for name := range js.StreamNames() {
		if err := js.DeleteStream(name); err != nil {
			t.Error(err)
		}
	}
}

func TestAddStream(t *testing.T) {
	info, err := js.AddStream(&nats.StreamConfig{
		Name:      "development",
		Subjects:  []string{"development"},
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestPublish(t *testing.T) {
	if _, err := js.Publish("development.message", []byte("hello")); err != nil {
		t.Error(err)
	}
}

func TestStreamInfo(t *testing.T) {
	v, err := js.StreamInfo("development")
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestDeleteStream(t *testing.T) {
	if err := js.DeleteStream("development"); err != nil {
		t.Error(err)
	}
}
