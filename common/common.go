package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
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
	STMP `yaml:"stmp"`
	NATS `yaml:"nats"`

	REDIS    string `yaml:"redis"`
	MONGO    string `yaml:"mongo"`
	MYSQL    string `yaml:"mysql"`
	POSTGRES string `yaml:"postgres"`
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

type NATS struct {
	Url  string `yaml:"url"`
	NKey string `yaml:"nkey"`
}

type Order struct {
	ID          uint64  `bun:"id,pk,autoincrement" faker:"-"`
	No          string  `bun:"type:varchar" faker:"cc_number"`
	Name        string  `bun:"type:varchar" faker:"name"`
	Description string  `bun:"type:text" faker:"paragraph"`
	Account     string  `bun:"type:varchar" faker:"username"`
	Customer    string  `bun:"type:varchar" faker:"name"`
	Email       string  `bun:"type:varchar" faker:"email"`
	Phone       string  `bun:"type:varchar" faker:"phone_number"`
	Address     string  `bun:"type:varchar" faker:"sentence"`
	Price       float64 `bun:"type:decimal" faker:"amount"`
	CreateTime  time.Time
	UpdateTime  time.Time
}
