package factory

import (
	"bytes"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
	"testing"
	"time"
)

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
		Code: "123456",
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
