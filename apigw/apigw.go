package apigw

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic/decoder"
	"net/http"
	"net/url"
	"time"
)

type Option struct {
	Source    string
	SecretId  string
	SecretKey string
}

type KeyAuthorizationResult struct {
	Datetime string
	Content  string
}

func KeyAuthorization(opt Option) (r *KeyAuthorizationResult, err error) {
	r = new(KeyAuthorizationResult)
	timeLocation, _ := time.LoadLocation("Etc/GMT")
	r.Datetime = time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", r.Datetime, opt.Source)

	mac := hmac.New(sha1.New, []byte(opt.SecretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	r.Content = fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		opt.SecretId, sign)

	return
}

func GetCity(ip string, opt Option) (data map[string]interface{}, err error) {
	opt.Source = "market"
	var r *KeyAuthorizationResult
	if r, err = KeyAuthorization(opt); err != nil {
		return
	}

	baseUrl, _ := url.Parse("https://service-84xhy5xe-1304874079.gz.apigw.tencentcs.com/release/v4")
	path := "/ip/city/query"
	u := baseUrl.JoinPath(path)
	query := u.Query()
	query.Add("ip", ip)
	query.Encode()
	u.RawQuery = query.Encode()

	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	req.Header.Set("X-Source", opt.Source)
	req.Header.Set("X-Date", r.Datetime)
	req.Header.Set("Authorization", r.Content)

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}
	if err = decoder.NewStreamDecoder(res.Body).Decode(&data); err != nil {
		return
	}
	return
}
