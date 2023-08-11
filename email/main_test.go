package email

import (
	"bytes"
	"crypto/tls"
	"development/common"
	"github.com/jordan-wright/email"
	"github.com/thoas/go-funk"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"testing"
	"time"
)

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

type Content struct {
	Name string
	User string
	Code string
	Year int
}

func TestSendVerifyCode(t *testing.T) {
	tpl, err := template.ParseFiles("./templates/verify_code.gohtml")
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, Content{
		Name: "WEPLANX",
		User: "Kain",
		Code: funk.RandomString(6, []rune("0123456789")),
		Year: time.Now().Year(),
	}); err != nil {
		t.Error(err)
	}
	e := &email.Email{
		To:      []string{"zhangtqx@vip.qq.com"},
		From:    "WEPLANX <weplanx@kainonly.com>",
		Subject: "验证",
		HTML:    buf.Bytes(),
	}
	if err = e.SendWithTLS(
		values.STMP.Addr,
		smtp.PlainAuth(
			values.STMP.Identity,
			values.STMP.Username,
			values.STMP.Password,
			values.STMP.Host,
		),
		&tls.Config{ServerName: values.STMP.Host},
	); err != nil {
		t.Error(err)
	}
}
