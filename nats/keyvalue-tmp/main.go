package main

import (
	"development/common"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
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
	kv, err := js.KeyValue("tmp")
	if err != nil {
		log.Fatal(err)
	}

	heart := time.NewTicker(time.Second)
	counter := 0
	go func() {
		for range heart.C {
			fmt.Println(counter)
			if counter == 10 {
				heart.Stop()
			}
			if _, err := kv.PutString("weplanx", "x1"); err != nil {
				log.Fatalln(err)
			}
			counter++
		}
	}()

	time.Sleep(time.Second)

	audience := time.NewTicker(time.Second)
	go func() {
		for range audience.C {
			entry, err := kv.Get("weplanx")
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(entry.Value()))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
