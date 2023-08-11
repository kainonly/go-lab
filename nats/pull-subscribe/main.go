package main

import (
	"development/common"
	"github.com/nats-io/nats.go"
	"log"
)

var nc *nats.Conn
var js nats.JetStreamContext

func main() {
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
