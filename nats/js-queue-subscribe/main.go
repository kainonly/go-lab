package main

import (
	"development/common"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
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

	if _, err := js.QueueSubscribe("development.message", "development:message", func(msg *nats.Msg) {
		log.Println("n1", string(msg.Data))
		msg.Term()
	}, nats.ManualAck()); err != nil {
		log.Fatalln(err)
	}
	if _, err := js.QueueSubscribe("development.message", "development:message", func(msg *nats.Msg) {
		log.Println("n2", string(msg.Data))
		msg.Ack()
	}, nats.ManualAck()); err != nil {
		log.Fatalln(err)
	}
	if _, err := js.QueueSubscribe("development.message", "development:message", func(msg *nats.Msg) {
		log.Println("n3", string(msg.Data))
		msg.Ack()
	}, nats.ManualAck()); err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
