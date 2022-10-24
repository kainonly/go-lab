package nats

import (
	"github.com/bytedance/sonic"
	"github.com/imdario/mergo"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

func TestCreateObjectStore(t *testing.T) {
	_, err := js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: "development",
	})
	if err != nil {
		t.Error(err)
	}
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
	_, err = o.PutBytes("values", b)
	if err != nil {
		t.Error(err)
	}
}

func TestGetObject(t *testing.T) {
	o, err := js.ObjectStore("development")
	if err != nil {
		t.Error(err)
	}
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

func TestMergeValues(t *testing.T) {
	dto := Values{
		UserLoginFailedTimes: 6,
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
	err := mergo.Merge(&dto, values)
	if err != nil {
		t.Error(err)
	}
	t.Log(dto)
}

func TestListObject(t *testing.T) {
	o, err := js.ObjectStore("development")
	if err != nil {
		return
	}
	obs, err := o.List()
	if err != nil {
		return
	}

	for _, x := range obs {
		var b []byte
		if b, err = o.GetBytes(x.Name); err != nil {
			return
		}
		t.Log(string(b))
	}
}

func TestDeleteObject(t *testing.T) {
	o, err := js.ObjectStore("development")
	if err != nil {
		t.Error(err)
	}
	if err = o.Delete("values"); err != nil {
		t.Error(err)
	}
}

func TestDeleteBucket(t *testing.T) {
	err := js.DeleteObjectStore("development")
	if err != nil {
		t.Error(err)
	}
}
