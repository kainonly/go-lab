package mock

import (
	"github.com/go-faker/faker/v4"
	"testing"
)

type Order struct {
	No          string  `faker:"cc_number"`
	Name        string  `faker:"name"`
	Description string  `faker:"paragraph"`
	Account     string  `faker:"username"`
	Customer    string  `faker:"name"`
	Email       string  `faker:"email"`
	Phone       string  `faker:"phone_number"`
	Address     string  `faker:"sentence"`
	Price       float64 `faker:"amount"`
	Time        string  `faker:"timestamp"`
}

func TestOrders(t *testing.T) {
	orders := make([]Order, 5)
	for i := 0; i < 5; i++ {
		if err := faker.FakeData(&orders[i]); err != nil {
			t.Error(err)
		}
	}
	t.Log(orders)
}
