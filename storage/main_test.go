package main

import (
	"development/common"
	"github.com/go-resty/resty/v2"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

var values *common.Values
var client *cos.Client
var httpclient *resty.Client

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	u, _ := url.Parse(values.COS.Url)
	b := &cos.BaseURL{BucketURL: u}
	client = cos.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  values.COS.AccessKeyID,
			SecretKey: values.COS.AccessKeySecret,
		},
	})
	httpclient = resty.New()
	os.Exit(m.Run())
}
