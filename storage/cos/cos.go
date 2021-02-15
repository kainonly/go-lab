package cos

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var (
	DEF = option{}
)

type option struct {
	Bucket    string
	Region    string
	SecretID  string
	SecretKey string
}

func Client() *cos.Client {
	u, _ := url.Parse("https://" + DEF.Bucket + ".cos." + DEF.Region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  DEF.SecretID,
			SecretKey: DEF.SecretKey,
		},
	})
}

func Put(key string) {

}

func Delete() {

}

func GeneratePostPresigned() {

}
