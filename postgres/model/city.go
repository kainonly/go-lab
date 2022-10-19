package model

type City struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar"`
	CountryCode string `gorm:"type:varchar;index"`
	StateCode   string `gorm:"type:varchar;index"`
	Latitude    float64
	Longitude   float64
}

type KVCity struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"`
	Value CityValue `gorm:"type:jsonb;serializer:json"`
}

type CityValue struct {
	Name        string  `json:"name"`
	CountryCode string  `json:"country_code"`
	StateCode   string  `json:"state_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
