package main

import (
	"github.com/nats-io/nats.go"
	"go-lab/common"
	"log"
	"time"
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
	o, err := js.ObjectStore("development")
	if err != nil {
		log.Fatal(err)
	}
	watch, err := o.Watch()
	if err != nil {
		log.Fatal(err)
	}
	cur := time.Now()
	for x := range watch.Updates() {
		if x == nil || x.ModTime.Unix() < cur.Unix() {
			continue
		}
		log.Println(x)
		log.Println(x.ModTime)
		log.Println(x.Deleted)
	}
}
