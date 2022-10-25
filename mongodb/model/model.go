package model

import (
	"time"
)

type Record struct {
	Time  time.Time `bson:"time"`
	Email string    `bson:"email" faker:"email"`
	Name  string    `bson:"name" faker:"name"`
	Ip    string    `bson:"ip" faker:"ipv4"`
}
