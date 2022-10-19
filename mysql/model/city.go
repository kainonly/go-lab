package model

type City struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255)"`
	CountryCode string `gorm:"type:varchar(50);index"`
	StateCode   string `gorm:"type:varchar(50);index"`
	Latitude    float64
	Longitude   float64
}

type KVCity struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"`
	Value CityValue `gorm:"type:json;serializer:json"`
}

type CityValue struct {
	Name        string  `json:"name"`
	CountryCode string  `json:"country_code"`
	StateCode   string  `json:"state_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
