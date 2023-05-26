package factory

import (
	"development/common"
	"log"
	"os"
	"testing"
)

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}
