package nats

import (
	"development/nats/common"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"testing"
)

var nc *nats.Conn
var js nats.JetStreamContext

func TestMain(m *testing.M) {
	os.Chdir("../")
	var err error
	if nc, err = common.Create("./config.yml"); err != nil {
		log.Fatalln(err)
	}
	if js, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
	); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestVersion(t *testing.T) {

}
