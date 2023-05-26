package model

type Model struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt int64
	UpdatedAt int64
}
