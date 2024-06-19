package common

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
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

func UseNats(values *Values) (nc *nats.Conn, err error) {
	var kp nkeys.KeyPair
	if kp, err = nkeys.FromSeed([]byte(values.NATS.NKey)); err != nil {
		panic(err)
	}

	defer kp.Wipe()

	var pub string
	if pub, err = kp.PublicKey(); err != nil {
		return
	}

	if nc, err = nats.Connect(
		values.NATS.Url,
		nats.Nkey(pub, func(nonce []byte) ([]byte, error) {
			sig, _ := kp.Sign(nonce)
			return sig, nil
		}),
	); err != nil {
		return
	}
	return
}
