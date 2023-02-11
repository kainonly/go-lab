package model

type Order struct {
	ID    uint64  `faker:"-"`
	Price float64 `faker:"amount"`
}
