package db

import (
	"context"
	"database/sql"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golab/common"
	"os"
	"testing"
)

var values *common.Values
var db *bun.DB
var mgo *mongo.Client
var mdb *mongo.Database
var rdb *redis.Client
var es *elasticsearch.Client
var influx influxdb2.Client
var msgId string

type M = map[string]interface{}

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		panic(err)
	}
	sqldb, err := sql.Open("mysql", values.MYSQL)
	if err != nil {
		panic(err)
	}
	db = bun.NewDB(sqldb, mysqldialect.New())
	if mgo, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(values.MONGO),
	); err != nil {
		panic(err)
	}
	option := options.Database().
		SetWriteConcern(writeconcern.Majority())
	mdb = mgo.Database("example", option)
	opts, err := redis.ParseURL(values.REDIS)
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opts)
	if es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: values.ELASTIC.Hosts,
		Username:  values.ELASTIC.Username,
		Password:  values.ELASTIC.Password,
	}); err != nil {
		panic(err)
	}
	influx = influxdb2.NewClient(
		values.INFLUX.Url,
		values.INFLUX.Token,
	)
	msgId = uuid.New().String()
	os.Exit(m.Run())
}

type Order struct {
	No          string  `bun:"type:varchar(20)" json:"no" bson:"no" faker:"cc_number"`
	Name        string  `bun:"type:varchar(50)" json:"name" bson:"name" faker:"name"`
	Description string  `bun:"type:varchar(1000)" json:"description" bson:"description" faker:"paragraph"`
	Account     string  `bun:"type:varchar(50)" json:"account" bson:"account" faker:"username"`
	Customer    string  `bun:"type:varchar(50)" json:"customer" bson:"customer" faker:"name"`
	Email       string  `bun:"type:varchar(50)" json:"email" bson:"email" faker:"email"`
	Phone       string  `bun:"type:varchar(20)" json:"phone" bson:"phone" faker:"phone_number"`
	Address     string  `bun:"type:varchar(255)" json:"address" bson:"address" faker:"sentence"`
	Price       float64 `bun:"type:decimal" json:"price" bson:"price" faker:"amount"`
}

type IOrder struct {
	bun.BaseModel `bun:"table:order"`
	ID            uint64 `bun:"id,pk,autoincrement" faker:"-"`
	*Order
}
