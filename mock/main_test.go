package mock

import (
	"github.com/bxcodec/faker/v3"
	"testing"
)

type Order struct {
	Customer string  `faker:"name"`
	Phone    string  `faker:"phone_number"`
	Cost     float64 `faker:"amount"`
	Time     string  `faker:"timestamp"`
}

func TestOrders(t *testing.T) {
	orders := make([]Order, 10)
	for i := 0; i < 10; i++ {
		if err := faker.FakeData(&orders[i]); err != nil {
			t.Error(err)
		}
	}
	t.Log(orders)
}
