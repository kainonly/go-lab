package vhost

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"text/template"
)

type ObjectVar struct {
	Domain     string
	BucketName string
}

func TestPutObject(t *testing.T) {
	tpl, err := template.ParseFiles("./object.txt")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, ObjectVar{
		Domain:     "console.kainonly.com",
		BucketName: "console",
	})
	assert.NoError(t, err)
	err = os.WriteFile("console.kainonly.com.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}

type CommonVar struct {
	Domain string
	Path   string
}

func TestPutStatic(t *testing.T) {
	tpl, err := template.ParseFiles("./static.txt")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, CommonVar{
		Domain: "www.kainonly.com",
		Path:   "/website/www",
	})
	err = os.WriteFile("www.kainonly.com.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}

func TestPutThinkPHP(t *testing.T) {
	tpl, err := template.ParseFiles("./thinkphp.txt")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, CommonVar{
		Domain: "wx.kainonly.com",
		Path:   "/website/wechat/public",
	})
	err = os.WriteFile("wx.kainonly.com.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}
