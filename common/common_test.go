package common

import (
	"net"
	"testing"
)

func TestNetIp(t *testing.T) {
	ip := net.ParseIP("119.41.34.152")
	t.Log(ip.To4())
}
