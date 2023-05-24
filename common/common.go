package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadValues(path string) (values *Values, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("静态配置不存在，请检查路径 [%s]", path)
	}
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &values); err != nil {
		return
	}
	return
}

type Values struct {
	CLS    `yaml:"cls"`
	STMP   `yaml:"stmp"`
	INFLUX `yaml:"influx"`
	PULSAR `yaml:"pulsar"`
	NATS   `yaml:"nats"`
	COS    `yaml:"cos"`

	REDIS      string `yaml:"redis"`
	MONGO      string `yaml:"mongo"`
	MYSQL      string `yaml:"mysql"`
	POSTGRES   string `yaml:"postgres"`
	POSTGREX   string `yaml:"postgrex"`
	CLICKHOUSE string `yaml:"clickhouse"`
	TencentBee string `yaml:"tencent_bee"`

	ELASTICSEARCH `yaml:"elasticsearch"`
	IPSERVICE     `yaml:"ip_service"`
}

type CLS struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	TopicId         string `yaml:"topic_id"`
}

type STMP struct {
	Addr     string `yaml:"addr"`
	Host     string `yaml:"host"`
	Identity string `yaml:"identity"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type INFLUX struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type PULSAR struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
	Topic string `yaml:"topic"`
}

type NATS struct {
	Url  string `yaml:"url"`
	NKey string `yaml:"nkey"`
}

type COS struct {
	Url             string `yaml:"url"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

type ELASTICSEARCH struct {
	Hosts    string `yaml:"hosts"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type IPSERVICE struct {
	SECRETID  string `yaml:"secretid"`
	SECRETKEY string `yaml:"secretkey"`
}

type Storage struct {
	Minio       StorageDrive `yaml:"minio"`
	Aliyun      StorageDrive `yaml:"aliyun"`
	Tencent     StorageDrive `yaml:"tencent"`
	Huaweicloud StorageDrive `yaml:"huaweicloud"`
	S3          StorageDrive `yaml:"s3"`
}

type StorageDrive struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	Cdn             string `yaml:"cdn,omitempty"`
}
