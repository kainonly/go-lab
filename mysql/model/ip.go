package model

type Ipv4 struct {
	ID       uint64 `gorm:"primaryKey,autoIncrement"`
	Start    uint64 `gorm:"index:idx_range"`
	End      uint64 `gorm:"index:idx_range"`
	Country  string `gorm:"type:varchar(50)"`
	Province string `gorm:"type:varchar(50)"`
	City     string `gorm:"type:varchar(50)"`
	ISP      string `gorm:"type:varchar(50)"`
}
