package storage

import (
	"development/common"
	"log"
	"os"
	"testing"
)

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	path := "../config/config.yml"
	if values, err = common.LoadValues(path); err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}
