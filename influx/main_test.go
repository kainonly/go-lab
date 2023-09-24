package influx

import (
	"context"
	"development/common"
	"fmt"
	"github.com/gookit/goutil/strutil"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var client influxdb2.Client
var msgId string

type M = map[string]interface{}

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		log.Fatalln(err)
	}

	client = influxdb2.NewClient(
		values.INFLUX.Url,
		values.INFLUX.Token,
	)

	msgId = strutil.MicroTimeHexID()
	os.Exit(m.Run())
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

func TestWriteJobs(t *testing.T) {
	now := time.Now()
	defer client.Close()
	api := client.WriteAPIBlocking("weplanx", "development")
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		tags := map[string]string{
			"key":   "65047ebfaa2334fff2f2a49c",
			"index": "0",
		}
		fields := M{
			"mode":   "HTTP",
			"url":    "https://whoami.kainonly.com/api",
			"header": nil,
			"body":   nil,
			"response": M{
				"status": 200,
				"body":   `{"hostname":"whoami-75d55b64f6-b5sgs","ip":["127.0.0.1","::1","10.42.1.108","fe80::b4a9:e6ff:fe78:e1e9"],"headers":{"Accept-Encoding":["gzip"],"Content-Length":["4"],"Content-Type":["application/json; charset=utf-8"],"User-Agent":["req/v3 (https://github.com/imroc/req)"],"X-Forwarded-For":["119.41.38.198"],"X-Forwarded-Host":["whoami.kainonly.com"],"X-Forwarded-Port":["443"],"X-Forwarded-Proto":["https"],"X-Forwarded-Server":["VM-12-17-opencloudos"],"X-Real-Ip":["119.41.38.198"]},"url":"/api","host":"whoami.kainonly.com","method":"POST","remoteAddr":"10.42.1.1:53340"}`,
			},
		}
		point := write.NewPoint("jobs", tags, fields, now.Add(-time.Second*5*time.Duration(i)))
		err := api.WritePoint(ctx, point)
		assert.NoError(t, err)
	}
}

func TestQueryJobsCount(t *testing.T) {
	defer client.Close()
	api := client.QueryAPI("weplanx")
	result, err := api.Query(context.TODO(), fmt.Sprintf(`
		from(bucket: "development")
		  |> range(start: -24h)
		  |> filter(fn: (r) => r["_measurement"] == "jobs")
		  |> filter(fn: (r) => r["key"] == "65047ebfaa2334fff2f2a49c")
		  |> filter(fn: (r) => r["index"] == "0")
		  |> count()
	`))
	assert.NoError(t, err)
	for result.Next() {
		t.Log(result.Record().Value())
	}
}

func TestQueryJobs(t *testing.T) {
	defer client.Close()
	api := client.QueryAPI("weplanx")
	result, err := api.Query(context.TODO(), fmt.Sprintf(`
		from(bucket: "development")
		  |> range(start: -1h)
		  |> filter(fn: (r) => r["_measurement"] == "jobs")
		  |> filter(fn: (r) => r["key"] == "65047ebfaa2334fff2f2a49c")
		  |> filter(fn: (r) => r["index"] == "0")
		  |> range(start: -15m)
	`))
	assert.NoError(t, err)
	for result.Next() {
		t.Log(result.Record().String())
	}
}
