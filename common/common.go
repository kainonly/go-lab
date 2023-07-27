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
	REDIS    string `yaml:"redis"`
	MONGO    string `yaml:"mongo"`
	MYSQL    string `yaml:"mysql"`
	POSTGRES string `yaml:"postgres"`

	STMP struct {
		Addr     string `yaml:"addr"`
		Host     string `yaml:"host"`
		Identity string `yaml:"identity"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"stmp"`

	NATS struct {
		Url  string `yaml:"url"`
		NKey string `yaml:"nkey"`
	} `yaml:"nats"`

	Apigw struct {
		Ip struct {
			SecretID  string `yaml:"secret_id"`
			SecretKey string `yaml:"secret_key"`
		} `yaml:"ip"`
	} `yaml:"apigw"`
}
