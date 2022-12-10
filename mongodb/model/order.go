package model

type Order struct {
	No          string  `bson:"no" faker:"cc_number"`
	Name        string  `bson:"name" faker:"name"`
	Description string  `bson:"description" faker:"paragraph"`
	Account     string  `bson:"account" faker:"username"`
	Customer    string  `bson:"customer" faker:"name"`
	Email       string  `bson:"email" faker:"email"`
	Phone       string  `bson:"phone" faker:"phone_number"`
	Address     string  `bson:"address" faker:"sentence"`
	Price       float64 `bson:"price" faker:"amount"`
}
