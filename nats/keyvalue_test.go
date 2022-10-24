package nats

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

func TestCreateKeyValue(t *testing.T) {
	_, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "development",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestKeyValueStores(t *testing.T) {
	for x := range js.KeyValueStores() {
		t.Log(x)
	}
}

func TestKeyValue(t *testing.T) {
	kv, err := js.KeyValue("development")
	if err != nil {
		t.Error(err)
	}
	var keys []string
	if keys, err = kv.Keys(); err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			keys = make([]string, 0)
		} else {
			t.Error(err)
		}
	}
	t.Log(keys)
}

func TestKeyValuePut(t *testing.T) {
	kv, err := js.KeyValue("development")
	if err != nil {
		t.Error(err)
	}
	values := Values{
		UserSessionExpire:    time.Hour,
		UserLoginFailedTimes: 5,
		UserLockTime:         time.Minute * 15,
		IpLoginFailedTimes:   5,
		IpWhitelist:          []string{},
		IpBlacklist:          []string{},
		PasswordStrength:     1,
	}
	var b []byte
	if b, err = sonic.Marshal(values); err != nil {
		t.Error(err)
	}
	r, err := kv.Put("values", b)
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestKeyValueGet(t *testing.T) {
	kv, err := js.KeyValue("development")
	if err != nil {
		t.Error(err)
	}
	entry, err := kv.Get("values")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(entry.Value()))
}
