package main

import (
	"github.com/nats-io/nats.go"
	"golab/common"
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
	kv, err := js.KeyValue("development")
	if err != nil {
		log.Fatal(err)
	}
	watch, err := kv.WatchAll()
	if err != nil {
		log.Fatal(err)
	}
	cur := time.Now()
	for entry := range watch.Updates() {
		if entry == nil || entry.Created().Unix() < cur.Unix() {
			continue
		}
		log.Println(entry.Operation().String())
		log.Println(entry.Key())
		log.Println(string(entry.Value()))
		log.Println(entry.Revision())
	}
}
