package apigw

import (
	"github.com/stretchr/testify/assert"
	"golab/common"
	"log"
	"os"
	"testing"
)

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestGetCity(t *testing.T) {
	data, err := GetCity("119.41.34.152", Option{
		SecretId:  values.Apigw.Ip.SecretID,
		SecretKey: values.Apigw.Ip.SecretKey,
	})
	assert.NoError(t, err)
	t.Log(data)
}
