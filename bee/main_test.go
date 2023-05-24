package bee

import (
	"development/common"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var values *common.Values
var client *resty.Client

func baseURL(path string) string {
	return fmt.Sprintf(
		`https://raw.githubusercontent.com/dr5hn/countries-states-cities-database/master/%s`,
		path,
	)
}

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

func TestTags(t *testing.T) {
	r, err := client.R().
		SetQueryParam("per_page", "100").
		Get("/projects/101042/repository/tags")
	assert.NoError(t, err)
	lists := make([]map[string]interface{}, 0)
	err = sonic.Unmarshal(r.Body(), &lists)
	assert.NoError(t, err)
}
