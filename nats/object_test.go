package nats

import (
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateObjectStore(t *testing.T) {
	_, err := js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: "development",
	})
	assert.NoError(t, err)
}

type Values struct {
	UserSessionExpire    time.Duration `json:"user_session_expire"`
	UserLoginFailedTimes int           `json:"user_login_failed_times"`
	UserLockTime         time.Duration `json:"user_lock_time"`
	IpLoginFailedTimes   int           `json:"ip_login_failed_times"`
	IpWhitelist          []string      `json:"ip_whitelist"`
	IpBlacklist          []string      `json:"ip_blacklist"`
	PasswordStrength     int           `json:"password_strength"`
}

func TestPutObject(t *testing.T) {
	o, err := js.ObjectStore("development")
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
	b, err = sonic.Marshal(values)
	assert.NoError(t, err)
	_, err = o.PutBytes("values", b)
	assert.NoError(t, err)
}

func TestGetObject(t *testing.T) {
	o, err := js.ObjectStore("development")
	assert.NoError(t, err)
	b, err := o.GetBytes("values")
	if err != nil {
		if err == nats.ErrObjectNotFound {
			t.Log(nats.ErrObjectNotFound)
		} else {
			t.Error(err)
		}
	}
	var values Values
	if err = sonic.Unmarshal(b, &values); err != nil {
		t.Error(err)
	}
	t.Log(values)
}

func TestListObject(t *testing.T) {
	o, err := js.ObjectStore("development")
	assert.NoError(t, err)
	obs, err := o.List()
	assert.NoError(t, err)

	for _, x := range obs {
		var b []byte
		b, err = o.GetBytes(x.Name)
		assert.NoError(t, err)
		t.Log(string(b))
	}
}

func TestDeleteObject(t *testing.T) {
	o, err := js.ObjectStore("development")
	assert.NoError(t, err)
	err = o.Delete("values")
	assert.NoError(t, err)
}

func TestDeleteBucket(t *testing.T) {
	err := js.DeleteObjectStore("development")
	assert.NoError(t, err)
}
