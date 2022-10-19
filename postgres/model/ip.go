package model

type Ipv4 struct {
	ID       uint64 `gorm:"primaryKey,autoIncrement"`
	Start    uint64 `gorm:"index:idx_range"`
	End      uint64 `gorm:"index:idx_range"`
	Country  string `gorm:"type:varchar"`
	Province string `gorm:"type:varchar"`
	City     string `gorm:"type:varchar"`
	ISP      string `gorm:"type:varchar"`
}
