package model

type OrderXL struct {
	ID          uint64  `gorm:"primaryKey,autoIncrement" faker:"-"`
	No          string  `gorm:"type:varchar(255)" faker:"cc_number"`
	Name        string  `gorm:"type:varchar(255)" faker:"name"`
	Description string  `gorm:"type:text" faker:"paragraph"`
	Account     string  `gorm:"type:varchar(255)" faker:"username"`
	Customer    string  `gorm:"type:varchar(255)" faker:"name"`
	Email       string  `gorm:"type:varchar(255)" faker:"email"`
	Phone       string  `gorm:"type:varchar(255)" faker:"phone_number"`
	Address     string  `gorm:"type:varchar(255)" faker:"sentence"`
	Price       float64 `gorm:"type:float" faker:"amount"`
}
