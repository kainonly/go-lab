package dbcomparison

import (
	"bytes"
	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"github.com/go-faker/faker/v4"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestEsInfo(t *testing.T) {
	resp, err := es.Info()
	assert.NoError(t, err)
	var result M
	err = decoder.NewStreamDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	t.Log(result)
}

func TestEsCreateIndex(t *testing.T) {
	resp, err := es.Indices.Create("orders")
	assert.NoError(t, err)
	var result M
	err = decoder.NewStreamDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	t.Log(result)
}

func TestEsDeleteIndex(t *testing.T) {
	resp, err := es.Indices.Delete([]string{"orders"})
	assert.NoError(t, err)
	var result M
	err = decoder.NewStreamDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	t.Log(result)
}

func TestEsIndex(t *testing.T) {
	var order Order
	err := faker.FakeData(&order)
	assert.NoError(t, err)
	var w = bytes.NewBuffer(nil)
	err = encoder.NewStreamEncoder(w).Encode(order)
	assert.NoError(t, err)
	resp, err := es.Index("orders", w)
	assert.NoError(t, err)
	var result M
	err = decoder.NewStreamDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	t.Log(result)
}

func TestEsBulk(t *testing.T) {
	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(100, func(i interface{}) {
		_, err := es.Bulk(i.(*bytes.Buffer))
		assert.NoError(t, err)
		wg.Done()
	})
	assert.NoError(t, err)
	defer p.Release()
	for n := 0; n < 100; n++ {
		wg.Add(1)
		w := bytes.NewBuffer(nil)
		stream := encoder.NewStreamEncoder(w)
		for i := 0; i < 10000; i++ {
			var order Order
			err = faker.FakeData(&order)
			assert.NoError(t, err)
			err = stream.Encode(M{"index": M{"_index": "orders"}})
			assert.NoError(t, err)
			err = stream.Encode(order)
			assert.NoError(t, err)
		}
		_ = p.Invoke(w)
	}
	wg.Wait()
}

func TestEsSearch(t *testing.T) {
	data := M{
		"query": M{
			"match_all": M{},
		},
	}
	var w = bytes.NewBuffer(nil)
	err := encoder.NewStreamEncoder(w).Encode(data)
	resp, err := es.Search(
		es.Search.WithIndex("orders"),
		es.Search.WithBody(w),
	)
	assert.NoError(t, err)
	var result M
	err = decoder.NewStreamDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	t.Log(result)
}
