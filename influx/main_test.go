package influx

import (
	"context"
	"development/common"
	"fmt"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var client influxdb2.Client
var msgId string

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	client = influxdb2.NewClient(
		values.INFLUX.Url,
		values.INFLUX.Token,
	)
	msgId = uuid.New().String()
	os.Exit(m.Run())
}

func TestWrite(t *testing.T) {
	defer client.Close()
	api := client.WriteAPI("weplanx", "development")
	p := influxdb2.NewPointWithMeasurement("beta").
		AddTag("msgId", msgId).
		AddField("num", 100).
		AddField("data", map[string]interface{}{
			"message": "ok",
			"sn": map[string]interface{}{
				"no": "123",
			},
		}).
		SetTime(time.Now())
	api.WritePoint(p)
	api.Flush()
}

func TestQuery(t *testing.T) {
	defer client.Close()
	api := client.QueryAPI("weplanx")
	result, err := api.Query(context.TODO(), fmt.Sprintf(`
		from(bucket: "development")
		  |> range(start: -1h)
		  |> filter(fn: (r) => r["_measurement"] == "beta")
		  |> filter(fn: (r) => r["msgId"] == "%s")
		  |> yield(name: "mean")
	`, msgId),
	)
	if err != nil {
		t.Error(err)
	}
	for result.Next() {
		t.Logf("value: %v\n", result.Record().Value())
	}
	if result.Err() != nil {
		t.Logf("query parsing error: %s\n", result.Err().Error())
	}
}
