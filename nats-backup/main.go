package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"golab/common"
	"log"
	"os"
	"os/exec"
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

	for name := range js.StreamNames() {
		cmd := exec.Command("nats", "stream", "backup", name, os.Getenv("BACKUP")+"\\"+name)
		output, _ := cmd.Output()
		fmt.Println(string(output))
	}
}
