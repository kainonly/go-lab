package mongodb

import (
	"context"
	"development/common"
	"development/mongodb/model"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/go-faker/faker/v4"
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

var values *common.Values
var client *mongo.Client
var db *mongo.Database

func TestMain(m *testing.M) {

	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	if client, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(values.MONGO),
	); err != nil {
		log.Fatalln(err)
	}

	option := options.Database().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	db = client.Database("xapi", option)
	os.Exit(m.Run())
}

func TestSort(t *testing.T) {
	ctx := context.TODO()
	option := options.Find().
		SetLimit(10).
		SetSkip(0).
		SetProjection(bson.M{"order_number": 1})

	cursor, err := db.Collection("orders").Find(ctx, bson.M{}, option)
	if err != nil {
		t.Error(err)
	}
	data := make([]map[string]interface{}, 0)
	if err = cursor.All(ctx, &data); err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestTransaction(t *testing.T) {
	ctx := context.TODO()

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	if err != nil {
		t.Error(err)
	}
	defer sess.EndSession(ctx)

	txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
	_, err = sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (result interface{}, err error) {
		if _, err = db.Collection("schema").DeleteOne(sessCtx, bson.M{"key": "role"}); err != nil {
			return
		}
		return nil, errors.New("test tx")
	}, txnOpts)
	if err != nil {
		t.Error(err)
	}
}

func TestTimeSeries(t *testing.T) {
	ctx := context.TODO()

	option := options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().SetTimeField("time"),
		)

	if err := db.CreateCollection(ctx, "history", option); err != nil {
		t.Error(err)
	}

	//docs := make([]interface{}, 100)
	//for i := 0; i < 100; i++ {
	//	var doc model.Record
	//	if err := faker.FakeData(&doc); err != nil {
	//		t.Error(err)
	//	}
	//	docs[i] = doc
	//}
	//
	//if _, err := db.Collection("history").InsertMany(ctx, docs); err != nil {
	//	t.Error(err)
	//}
}

func TestExistsTimeSeriesDb(t *testing.T) {
	ctx := context.TODO()

	colls, err := db.ListCollectionSpecifications(ctx,
		bson.M{
			"name": "history",
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(colls) != 0 {
		t.Log(colls[0].Type)
	}
}

func TestCreateUser(t *testing.T) {
	hash, err := argon2id.CreateHash("pass@VAN1234", &argon2id.Params{
		Memory:      65536,
		Iterations:  4,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
	if err != nil {
		t.Error(err)
	}

	ctx := context.TODO()

	if _, err := db.Collection("users").InsertOne(ctx, model.User{
		Username:   "weplanx",
		Password:   hash,
		Roles:      []primitive.ObjectID{},
		Email:      "zhangtqx@vip.qq.com",
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		t.Error(err)
	}

	if _, err := db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetName("idx_username").SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetName("idx_email"),
		},
	}); err != nil {
		t.Error(err)
	}
}

func TestCreateProjectsCollection(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetValidator(bson.D{
			{"$jsonSchema", bson.D{
				{"title", "projects"},
				{"required", bson.A{"_id", "name", "namespace", "status", "create_time", "update_time"}},
				{"properties", bson.D{
					{"_id", bson.M{"bsonType": "objectId"}},
					{"name", bson.M{"bsonType": "string"}},
					{"namespace", bson.M{"bsonType": "string"}},
					{"secret", bson.M{"bsonType": []string{"null", "string"}}},
					{"entry", bson.M{"bsonType": "array"}},
					{"expire_time", bson.M{"bsonType": []string{"null", "date"}}},
					{"labels", bson.M{"bsonType": "object"}},
					{"status", bson.M{"bsonType": "bool"}},
					{"create_time", bson.M{"bsonType": "date"}},
					{"update_time", bson.M{"bsonType": "date"}},
				}},
				{"additionalProperties", false},
			}},
		})
	if err := db.CreateCollection(ctx, "projects", option); err != nil {
		t.Error(err)
	}
	if _, err := db.Collection("projects").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"namespace", 1}},
			Options: options.Index().SetName("idx_namespace").SetUnique(true),
		},
	}); err != nil {
		t.Error(err)
	}
}

func TestCreateProject(t *testing.T) {
	ctx := context.TODO()
	if _, err := db.Collection("projects").InsertOne(ctx, model.Project{
		Name:      "默认项目",
		Namespace: "default",
		Entry:     []string{},
		Labels: map[string]string{
			"fixed": "true",
		},
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		t.Error(err)
	}
}

func TestMockOrder(t *testing.T) {
	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(100, func(i interface{}) {
		if _, err := db.Collection("orders").InsertMany(context.TODO(), i.([]interface{})); err != nil {
			t.Error(err)
		}
		wg.Done()
	})
	if err != nil {
		t.Error(err)
	}
	defer p.Release()
	for n := 0; n < 100; n++ {
		wg.Add(1)
		orders := make([]interface{}, 10000)
		for i := 0; i < 10000; i++ {
			var data model.Order
			if err := faker.FakeData(&data); err != nil {
				t.Error(err)
			}
			orders[i] = data
		}
		_ = p.Invoke(orders)
	}
	wg.Wait()
}

func TestMockOrderXL(t *testing.T) {
	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(100, func(i interface{}) {
		if _, err := db.Collection("orders_xl").InsertMany(context.TODO(), i.([]interface{})); err != nil {
			t.Error(err)
		}
		wg.Done()
	})
	if err != nil {
		t.Error(err)
	}
	defer p.Release()
	for n := 0; n < 1000; n++ {
		wg.Add(1)
		orders := make([]interface{}, 10000)
		for i := 0; i < 10000; i++ {
			var data model.Order
			if err := faker.FakeData(&data); err != nil {
				t.Error(err)
			}
			orders[i] = data
		}
		_ = p.Invoke(orders)
	}
	wg.Wait()
}

func TestMockDevTable(t *testing.T) {
	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(100, func(i interface{}) {
		if _, err := db.Collection("dev_table").InsertMany(context.TODO(), i.([]interface{})); err != nil {
			t.Error(err)
		}
		wg.Done()
	})
	if err != nil {
		t.Error(err)
	}
	defer p.Release()
	for n := 0; n < 100; n++ {
		wg.Add(1)
		orders := make([]interface{}, 10000)
		for i := 0; i < 10000; i++ {
			var data model.DevTable
			if err := faker.FakeData(&data); err != nil {
				t.Error(err)
			}
			data.CreateTime, _ = time.Parse(`2006-01-02 15:04:05`, faker.Timestamp())
			data.UpdateTime = data.CreateTime.Add(time.Hour * 24)
			orders[i] = data
		}
		_ = p.Invoke(orders)
	}

}
