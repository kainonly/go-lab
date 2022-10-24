package main

import (
	"development/nats/common"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	nc, err := common.Create("./config/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
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
