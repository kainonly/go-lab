package main

type Order struct {
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

type OrderA struct {
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

type OrderB struct {
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

type OrderC struct {
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

type Tag struct {
	ID   uint64 `gorm:"primaryKey,autoIncrement" faker:"-"`
	Name string `gorm:"type:varchar(255)" faker:"name"`
}

type TUser struct {
	ID    uint64 `gorm:"primaryKey,autoIncrement" faker:"-"`
	Name  string `gorm:"type:varchar(255)" faker:"name"`
	TagId uint64
}

type TOrder struct {
	ID      uint64  `gorm:"primaryKey,autoIncrement" faker:"-"`
	No      string  `gorm:"type:varchar(255)" faker:"cc_number"`
	Price   float64 `gorm:"type:float" faker:"amount"`
	TUserId uint64
}
