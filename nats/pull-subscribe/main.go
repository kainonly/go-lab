package main

import (
	"development/nats/common"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, _ := common.Create("./config/config.yml")
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
	sub, err := js.PullSubscribe("development.message", "development:message")
	if err != nil {
		panic(err)
	}
	msgs, err := sub.Fetch(10)
	if err != nil {
		panic(err)
	}
	log.Println(msgs)
}
