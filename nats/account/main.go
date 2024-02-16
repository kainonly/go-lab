package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go-lab/common"
	"os"
	"os/signal"
)

var nc *nats.Conn
var js nats.JetStreamContext

func UseNats(v *common.Values) (err error) {
	if nc, err = nats.Connect(
		v.NATS.Url,
		nats.UserInfo("test", "secret"),
	); err != nil {
		return
	}
	if js, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
	); err != nil {
		return
	}
	return
}

func main() {
	values, err := common.LoadValues("./config.yml")
	if err != nil {
		panic(err)
	}
	if err = UseNats(values); err != nil {
		panic(err)
	}

	if _, err = js.AddStream(&nats.StreamConfig{Name: "test", Subjects: []string{"test"}}); err != nil {
		panic(err)
	}

	if _, err = js.Subscribe("test", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	}); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
