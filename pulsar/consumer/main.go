package main

import (
	"development/common"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

func main() {
	var err error
	var values *common.Values
	if values, err = common.LoadValues(); err != nil {
		log.Fatalln(err)
	}
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		// 服务接入地址
		URL: values.PULSAR.Url,
		// 授权角色密钥
		Authentication:    pulsar.NewAuthenticationToken(values.PULSAR.Token),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}
	defer client.Close()

	go func() {
		// 使用客户端创建消费者
		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:            values.PULSAR.Topic,
			SubscriptionName: "development",
			Type:             pulsar.Shared,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer consumer.Close()
		for c := range consumer.Chan() {
			msg := c.Message
			fmt.Println(string(c.Message.Payload()))
			consumer.Ack(msg)
		}
	}()

	select {}
}
