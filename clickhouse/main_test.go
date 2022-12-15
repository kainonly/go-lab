package clickhouse

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
	if db, err = gorm.Open(clickhouse.Open(values.CLICKHOUSE)); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestQuery(t *testing.T) {
	var result map[string]interface{}
	if err := db.Raw(`select avg(price) from orders`).Scan(&result); err != nil {
		t.Error(err)
	}
	t.Log(result)

}
