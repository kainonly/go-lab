package ip

import (
	"crypto/hmac"
	"crypto/sha1"
	"development/common"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func urlencode(params map[string]string) string {
	var p = url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}

func TestQueryIp(t *testing.T) {
	source := "market"
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	datetime := time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", datetime, source)
	mac := hmac.New(sha1.New, []byte(values.IPSERVICE.SECRETID))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	auth := fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		values.IPSERVICE.SECRETKEY, sign)
	method := "GET"
	headers := map[string]string{
		"X-Source":      source,
		"X-Date":        datetime,
		"Authorization": auth,
	}
	queryParams := map[string]string{
		"ip": "115.183.152.152",
	}
	bodyParams := make(map[string]string)
	u := "https://service-8t32o0n9-1300755093.ap-beijing.apigateway.myqcloud.com/release/lifeservice/QueryIpAddr/query"
	if len(queryParams) > 0 {
		u = fmt.Sprintf("%s?%s", u, urlencode(queryParams))
	}

	bodyMethods := map[string]bool{"POST": true, "PUT": true, "PATCH": true}
	var body io.Reader = nil
	if bodyMethods[method] {
		body = strings.NewReader(urlencode(bodyParams))
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest(method, u, body)
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bodyBytes))
}
