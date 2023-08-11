package pulsar

import (
	"context"
	"development/common"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var client pulsar.Client

func TestMain(m *testing.M) {
	os.Chdir("../")
	var err error
	if values, err = common.LoadValues("./config.yml"); err != nil {
		log.Fatalln(err)
	}
	if client, err = pulsar.NewClient(pulsar.ClientOptions{
		// 服务接入地址
		URL: values.PULSAR.Url,
		// 授权角色密钥
		Authentication:    pulsar.NewAuthenticationToken(values.PULSAR.Token),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	}); err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}
	os.Exit(m.Run())
}

func TestCreateProducer(t *testing.T) {
	defer client.Close()
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: values.PULSAR.Topic,
	})
	if err != nil {
		t.Error(err)
	}
	defer producer.Close()
	// 发送消息
	_, err = producer.Send(context.TODO(), &pulsar.ProducerMessage{
		// 消息内容
		Payload: []byte(`{"msg":"hi"}`),
		//// 业务key
		//Key: "yourKey",
		//// 业务参数
		//Properties: map[string]string{"key": "value"},
	})
}
