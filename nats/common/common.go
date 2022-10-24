package common

import (
	"development/common"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"time"
)

func Create(path string) (nc *nats.Conn, err error) {
	var values *common.Values
	if values, err = common.LoadValues(path); err != nil {
		return
	}

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
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second),
		nats.ReconnectJitter(500*time.Millisecond, 2*time.Second),
		nats.Nkey(pub, func(nonce []byte) ([]byte, error) {
			sig, _ := kp.Sign(nonce)
			return sig, nil
		}),
	); err != nil {
		return
	}
	return
}
