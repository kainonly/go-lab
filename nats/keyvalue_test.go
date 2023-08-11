package nats

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateKeyValue(t *testing.T) {
	_, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "development",
	})
	assert.NoError(t, err)
}

func TestKeyValueStores(t *testing.T) {
	for x := range js.KeyValueStores() {
		t.Log(x)
	}
}

func TestKeyValue(t *testing.T) {
	kv, err := js.KeyValue("development")
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	t.Log(r)
}

func TestKeyValueGet(t *testing.T) {
	kv, err := js.KeyValue("development")
	assert.NoError(t, err)
	entry, err := kv.Get("values")
	assert.NoError(t, err)
	t.Log(string(entry.Value()))
}

func TestKeyValueDel(t *testing.T) {
	kv, err := js.KeyValue("development")
	assert.NoError(t, err)
	err = kv.Delete("values")
	assert.NoError(t, err)
}

func TestDeleteKeyValue(t *testing.T) {
	err := js.DeleteKeyValue("development")
	assert.NoError(t, err)
}
