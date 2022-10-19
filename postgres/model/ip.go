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

type IpValue struct {
	Start    uint64 `json:"start"`
	End      uint64 `json:"end"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

type KVIpv4 struct {
	ID    uint64  `gorm:"primaryKey,autoIncrement"`
	Value IpValue `gorm:"type:jsonb;serializer:json"`
}
