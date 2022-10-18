package loki

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/go-resty/resty/v2"
	"github.com/grafana/loki/pkg/logproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"testing"
	"time"
)

func TestGrpc(t *testing.T) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("development:9096", opts...)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	client := logproto.NewPusherClient(conn)
	res, err := client.Push(context.TODO(), &logproto.PushRequest{
		Streams: []logproto.Stream{
			{
				Labels: `{topic="test"}`,
				Entries: []logproto.Entry{
					{
						Timestamp: time.Now(),
						Line:      `level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg="compact blocks" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources="[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]" duration=2.897213221s`,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

type Stream struct {
	Stream map[string]interface{} `json:"stream"`
	Values [][]interface{}        `json:"values"`
}

type Data struct {
	Email string `json:"email" faker:"email"`
	Name  string `json:"name" faker:"name"`
	Ip    string `json:"ip" faker:"ipv4"`
}

func TestHTTP(t *testing.T) {
	streams := make([]Stream, 1)
	for i := 0; i < 1; i++ {
		var data Data
		if err := faker.FakeData(&data); err != nil {
			t.Error(err)
		}
		streams[i] = Stream{
			Stream: map[string]interface{}{
				"email": data.Email,
			},
			Values: [][]interface{}{
				{strconv.Itoa(int(time.Now().UnixNano())), data.Ip},
			},
		}
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]interface{}{"streams": streams}).
		Post("http://development:3100/loki/api/v1/push")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
