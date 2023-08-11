package factory

import (
	"bytes"
	"development/common"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"text/template"
)

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

type ObjectVar struct {
	Domain     string
	BucketName string
}

func TestPutObject(t *testing.T) {
	tpl, err := template.ParseFiles("./templates/object.txt")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, ObjectVar{
		Domain:     "console.kainonly.com",
		BucketName: "console",
	})
	assert.NoError(t, err)
	err = os.WriteFile("./dist/console.kainonly.com.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}

type CommonVar struct {
	Domain string
	Path   string
}

func TestPutStatic(t *testing.T) {
	tpl, err := template.ParseFiles("./templates/static.txt")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, CommonVar{
		Domain: "www.kainonly.com",
		Path:   "/website/www",
	})
	err = os.WriteFile("./dist/www.kainonly.com.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}

func TestPutThinkPHP(t *testing.T) {
	tpl, err := template.ParseFiles("./templates/thinkphp.txt")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, CommonVar{
		Domain: "wx.kainonly.com",
		Path:   "/website/wechat/public",
	})
	err = os.WriteFile("./dist/wx.kainonly.com.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}
