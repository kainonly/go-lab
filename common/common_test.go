package common

import (
	"github.com/gookit/goutil/strutil"
	"net"
	"testing"
)

func TestNetIp(t *testing.T) {
	ip := net.ParseIP("119.41.34.152")
	t.Log(ip.To4())

}

func TestRandom(t *testing.T) {
	v := strutil.RandomChars(32)
	t.Log(v)
}
