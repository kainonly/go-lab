package postgres

import (
	"development/common"
	"development/mysql/model"
	"github.com/alexedwards/argon2id"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var values *common.Values
var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues(); err != nil {
		log.Fatalln(err)
	}
	if db, err = gorm.Open(mysql.Open(values.MYSQL), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	if err := db.AutoMigrate(model.User{}); err != nil {
		t.Error(err)
	}
	hash, err := argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams)
	if err != nil {
		t.Error(err)
	}
	if err = db.Create(&model.User{
		Username: "weplanx",
		Password: hash,
		Email:    "zhangtqx@vip.qq.com",
	}).Error; err != nil {
		t.Error(err)
	}
}
