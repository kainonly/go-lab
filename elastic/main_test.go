package main

import (
	"bytes"
	"development/common"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

var values *common.Values
var es *elasticsearch.Client

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		log.Fatalln(err)
	}
	cfg := elasticsearch.Config{
		Addresses: strings.Split(values.ELASTICSEARCH.Hosts, ","),
		Username:  values.ELASTICSEARCH.Username,
		Password:  values.ELASTICSEARCH.Password,
	}
	if es, err = elasticsearch.NewClient(cfg); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestGetInfo(t *testing.T) {
	r, err := es.Info()
	assert.NoError(t, err)
	var data map[string]interface{}
	err = decoder.NewStreamDecoder(r.Body).Decode(&data)
	assert.NoError(t, err)
	t.Log(data)
}

func TestIndex(t *testing.T) {
	b, err := sonic.Marshal(map[string]interface{}{
		"msg": "hello",
	})
	assert.NoError(t, err)
	r, err := es.Index("test", bytes.NewReader(b))
	assert.NoError(t, err)
	var data map[string]interface{}
	err = decoder.NewStreamDecoder(r.Body).Decode(&data)
	assert.NoError(t, err)
	t.Log(data)
}

func TestMockOrder(t *testing.T) {
	for n := 0; n < 1000000; n++ {
		var data Order
		if err := faker.FakeData(&data); err != nil {
			t.Error(err)
		}
		b, err := sonic.Marshal(data)
		assert.NoError(t, err)
		_, err = es.Index("orders", bytes.NewReader(b))
		assert.NoError(t, err)
	}
}
