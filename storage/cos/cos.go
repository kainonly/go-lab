package cos

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-extra/str"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

var (
	DEF = Option{}
)

type Option struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
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

func Put(fileHeader *multipart.FileHeader) (fileName string, err error) {
	client := Client()
	fileName = time.Now().Format("20060102") +
		"/" + str.Uuid().String() +
		"." + path.Ext(fileHeader.Filename)
	var file multipart.File
	if file, err = fileHeader.Open(); err != nil {
		return
	}
	if _, err = client.Object.Put(
		context.Background(), fileName, file, nil,
	); err != nil {
		return
	}
	return
}

func Delete(keys []string) error {
	client := Client()
	objects := make([]cos.Object, len(keys))
	for index, value := range keys {
		objects[index] = cos.Object{
			Key: value,
		}
	}
	if _, _, err := client.Object.DeleteMulti(context.Background(), &cos.ObjectDeleteMultiOptions{
		Quiet:   true,
		Objects: objects,
	}); err != nil {
		return err
	}
	return nil
}

func GeneratePostPresigned(expired int64, conditions ...[]interface{}) (data map[string]interface{}, err error) {
	now := time.Now()
	keyTime := strconv.Itoa(int(now.Unix())) + ";" + strconv.Itoa(int(now.Unix()+expired))
	fileName := time.Now().Format("20060102") + "/" + str.Uuid().String()
	conditions = append(conditions, []interface{}{"bucket", DEF.Bucket})
	conditions = append(conditions, []interface{}{"starts-with", "$key", fileName})
	conditions = append(conditions, []interface{}{"q-sign-algorithm", "sha1"})
	conditions = append(conditions, []interface{}{"q-ak", DEF.SecretID})
	conditions = append(conditions, []interface{}{"q-sign-time", keyTime})
	policy := map[string]interface{}{
		"expiration": now.Add(time.Duration(expired) * time.Second).UTC().Format(time.RFC3339),
		"conditions": conditions,
	}
	var policyData []byte
	if policyData, err = jsoniter.Marshal(policy); err != nil {
		return
	}
	signKeyHash := hmac.New(sha1.New, []byte(DEF.SecretKey))
	signKeyHash.Write([]byte(keyTime))
	signKey := base64.StdEncoding.EncodeToString(signKeyHash.Sum(nil))
	stringToSignHash := sha1.New()
	stringToSignHash.Write(policyData)
	stringToSign := base64.StdEncoding.EncodeToString(stringToSignHash.Sum(nil))
	signatureHash := hmac.New(sha1.New, []byte(signKey))
	signatureHash.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(signatureHash.Sum(nil))
	data = map[string]interface{}{
		"filename": fileName,
		"option": map[string]interface{}{
			"ak":             DEF.SecretID,
			"policy":         base64.StdEncoding.EncodeToString(policyData),
			"key_time":       keyTime,
			"sign_algorithm": "sha1",
			"signature":      signature,
		},
	}
	return

}
