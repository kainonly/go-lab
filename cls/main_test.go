package cls

import (
	"development/common"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var producerInstance *cls.AsyncProducerClient
var topicId string

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues(); err != nil {
		log.Fatalln(err)
	}
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.Endpoint = values.CLS.Endpoint
	producerConfig.AccessKeyID = values.CLS.AccessKeyID
	producerConfig.AccessKeySecret = values.CLS.AccessKeySecret
	if producerInstance, err = cls.NewAsyncProducerClient(producerConfig); err != nil {
		log.Fatalln(err)
	}
	producerInstance.Start()
	topicId = values.CLS.TopicId
	os.Exit(m.Run())
}

type Callback struct {
	t *testing.T
}

func (x *Callback) Success(result *cls.Result) {
	attemptList := result.GetReservedAttempts()
	for _, attempt := range attemptList {
		x.t.Log(attempt)
	}
}

func (x *Callback) Fail(result *cls.Result) {
	x.t.Error(result.GetErrorMessage())
}

func TestSendLog(t *testing.T) {
	callBack := &Callback{t: t}
	message := cls.NewCLSLog(
		time.Now().Unix(),
		map[string]string{"content": "hello world| I'm from Beijing"},
	)
	if err := producerInstance.SendLog(topicId, message, callBack); err != nil {
		t.Error(err)
	}
	producerInstance.Close(60000)
}
