package postgres

import (
	"development/common"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var values *common.Values
var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	if db, err = gorm.Open(clickhouse.Open(values.CLICKHOUSE), &gorm.Config{}); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestAvg(t *testing.T) {
	var avg float64
	if err := db.Debug().Raw(`select avg(price) from orders`).Scan(&avg).Error; err != nil {
		t.Error(err)
	}
	t.Log(avg)
}
