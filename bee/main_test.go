package bee

import (
	"development/common"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var values *common.Values
var client *resty.Client

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	client = resty.New().
		SetBaseURL("https://git.code.tencent.com/api/v3").
		SetHeader("PRIVATE-TOKEN", values.TencentBee)
	os.Exit(m.Run())
}

func TestArchive(t *testing.T) {
	r, err := client.R().
		SetQueryParam("sha", "1.0.0").
		Get("/projects/100032/repository/archive")
	assert.NoError(t, err)
	err = os.WriteFile("1.0.0.zip", r.Body(), os.ModePerm)
	assert.NoError(t, err)
}
