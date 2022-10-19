package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func LoadValues() (values *Values, err error) {
	path := "../config/config.yml"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("静态配置不存在，请检查路径 [%s]", path)
	}
	var b []byte
	if b, err = ioutil.ReadFile(path); err != nil {
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

	REDIS string `yaml:"redis"`
	MONGO string `yaml:"mongo"`
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
