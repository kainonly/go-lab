package main

import (
	"development/common"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
)

var nc *nats.Conn

func main() {
	values, err := common.LoadValues("./config.yml")
	if err != nil {
		panic(err)
	}
	if nc, err = common.UseNats(values); err != nil {
		panic(err)
	}
	ns := []int{0, 0, 0}
	if err = run(ns); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println(ns)
}

func run(ns []int) (err error) {
	if _, err = nc.QueueSubscribe("development.message", "development:message", func(msg *nats.Msg) {
		fmt.Println("n1", string(msg.Data))
		ns[0]++
	}); err != nil {
		return
	}
	if _, err = nc.QueueSubscribe("development.message", "development:message", func(msg *nats.Msg) {
		fmt.Println("n2", string(msg.Data))
		ns[1]++
	}); err != nil {
		return
	}
	if _, err = nc.QueueSubscribe("development.message", "development:message", func(msg *nats.Msg) {
		fmt.Println("n3", string(msg.Data))
		ns[2]++
	}); err != nil {
		return
	}
	return
}
