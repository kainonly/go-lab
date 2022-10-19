package nats

import (
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"testing"
)

func TestConsumerNames(t *testing.T) {
	for name := range js.ConsumerNames("development") {
		t.Log(name)
	}
}

func TestConsumersInfo(t *testing.T) {
	//for v := range js.ConsumersInfo("development") {
	//	data, _ := json.Marshal(v)
	//	t.Log(string(data))
	//}
	for v := range js.Consumers("development") {
		data, _ := sonic.Marshal(v)
		t.Log(string(data))
	}
}

func TestAddConsumer(t *testing.T) {
	info, err := js.AddConsumer("development", &nats.ConsumerConfig{
		Durable: "DEV",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestConsumerInfo(t *testing.T) {
	info, err := js.ConsumerInfo("development", "DEV")
	if err != nil {
		t.Error(err)
	}
	data, _ := sonic.Marshal(info)
	t.Log(string(data))
}
