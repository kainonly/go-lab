package influx

import (
	"context"
	"development/common"
	"fmt"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stretchr/testify/assert"
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
	if values, err = common.LoadValues("./config.yml"); err != nil {
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

func TestMongoOpenConnections(t *testing.T) {
	defer client.Close()
	queryAPI := client.QueryAPI("weplanx")
	query := `option v = {timeRangeStart: -12h, timeRangeStop: now(), windowPeriod: 10000ms}
		from(bucket: "observability")
		|> range(start: v.timeRangeStart, stop: v.timeRangeStop)
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["hostname"] == "xmgo.kainonly.com")
		|> filter(fn: (r) => r["_field"] == "open_connections")
		|> derivative(unit: v.windowPeriod, nonNegative: true)
		|> yield(name: "nonnegative derivative")
	`
	result, err := queryAPI.Query(context.Background(), query)
	assert.NoError(t, err)
	for result.Next() {
		// Notice when group key has changed
		//if result.TableChanged() {
		//	fmt.Printf("table: %s\n", result.TableMetadata().String())
		//}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Value())
	}
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}
}

func TestCgoCalls(t *testing.T) {
	defer client.Close()
	queryAPI := client.QueryAPI("weplanx")
	query := `
		from(bucket: "development")
			|> range(start: -15m, stop: now())
			|> filter(fn: (r) => r["_measurement"] == "prometheus")
			|> filter(fn: (r) => r["_field"] == "process.runtime.go.cgo.calls")
			|> filter(fn: (r) => r["service.name"] == "weplanx")
			|> group(columns: ["service.name"], mode: "by")
			|> aggregateWindow(every: 5m, fn: mean, createEmpty: false)
			|> yield(name: "mean")
	`
	result, err := queryAPI.Query(context.Background(), query)
	assert.NoError(t, err)
	data := make([]interface{}, 0)
	for result.Next() {
		data = append(data, map[string]interface{}{
			"timestamp": result.Record().Time(),
			"value":     result.Record().Value(),
		})
	}
	assert.NoError(t, result.Err())
	t.Log(data)
}

func TestP99(t *testing.T) {
	defer client.Close()
	queryAPI := client.QueryAPI("weplanx")
	query := `
		import "experimental/aggregate"

		from(bucket: "development")
			|> range(start: -15m, stop: now())
			|> filter(fn: (r) => r["_measurement"] == "prometheus")
			|> filter(fn: (r) => r["service.name"] == "weplanx")
			|> filter(fn: (r) => r["_field"] == "http.server.duration_bucket")
			|> map(fn: (r) => ({r with le: float(v: r.le)}))
			|> aggregate.rate(every: 5m, groupColumns: ["le", "service.name", "http.method"])
			|> fill(value: 0.0)
			|> histogramQuantile(quantile: 0.99)
	`
	result, err := queryAPI.Query(context.Background(), query)
	assert.NoError(t, err)
	for result.Next() {
		fmt.Println(result.Record().String())
	}
	assert.NoError(t, result.Err())
}
