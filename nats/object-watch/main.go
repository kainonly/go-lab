package main

import (
	"development/nats/common"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	nc, _ := common.Create("./config/config.yml")
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
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
