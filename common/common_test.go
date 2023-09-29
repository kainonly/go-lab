package common

import (
	"math/rand"
	"net"
	"testing"
)

func TestNetIp(t *testing.T) {
	ip := net.ParseIP("119.41.34.152")
	t.Log(ip.To4())

}

func TestRandom(t *testing.T) {
	var data [][]int
	for i := 0; i < 10; i++ {
		n := rand.Intn(10)
		data = append(data, rand.Perm(10)[:n])
	}
	t.Log(data)
}
