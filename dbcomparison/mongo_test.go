package dbcomparison

import (
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"testing"
	"time"
)

func TestMgoSort(t *testing.T) {
	ctx := context.TODO()
	option := options.Find().
		SetLimit(10).
		SetSkip(0).
		SetProjection(bson.M{"order_number": 1})

	cursor, err := mdb.Collection("orders").Find(ctx, bson.M{}, option)
	assert.NoError(t, err)
	data := make([]map[string]interface{}, 0)
	err = cursor.All(ctx, &data)
	assert.NoError(t, err)
	t.Log(data)
}

func TestMgoTransaction(t *testing.T) {
	ctx := context.TODO()
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	session, err := mgo.StartSession(opts)
	assert.NoError(t, err)
	defer session.EndSession(ctx)

	txnOpts := options.Transaction().
		SetReadPreference(readpref.PrimaryPreferred())
	_, err = session.WithTransaction(ctx, func(sctx mongo.SessionContext) (_ interface{}, err error) {
		r, err := mdb.Collection("roles").InsertOne(sctx, bson.M{"name": "super"})
		assert.NoError(t, err)

		// TODO: 假设唯一索引出错
		_, err = mdb.Collection("users").InsertOne(sctx, bson.M{
			"name": "kain",
			"role": []primitive.ObjectID{r.InsertedID.(primitive.ObjectID)},
		})
		assert.NoError(t, err)
		return
	}, txnOpts)
	assert.NoError(t, err)
}

func TestMgoStartTransaction(t *testing.T) {
	ctx := context.TODO()
	opts := options.Session().
		SetDefaultReadConcern(readconcern.Majority())
	session, err := mgo.StartSession(opts)
	assert.NoError(t, err)
	defer session.EndSession(ctx)

	txnOpts := options.Transaction().
		SetReadPreference(readpref.PrimaryPreferred())
	err = session.StartTransaction(txnOpts)
	assert.NoError(t, err)

	err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) (err error) {
		r, err := mdb.Collection("roles").InsertOne(sctx, bson.M{"name": "super"})
		assert.NoError(t, err)
		if _, err = mdb.Collection("users").InsertOne(sctx, bson.M{
			"name": "kain",
			"role": []primitive.ObjectID{r.InsertedID.(primitive.ObjectID)},
		}); err != nil {
			session.AbortTransaction(sctx)
			return
		}
		session.CommitTransaction(sctx)
		return
	})
	assert.NoError(t, err)
}

func TestMgoTimeSeries(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().SetTimeField("time"),
		)
	err := mdb.CreateCollection(ctx, "history", option)
	assert.NoError(t, err)
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

func TestMgoExistsTimeSeriesDb(t *testing.T) {
	ctx := context.TODO()
	colls, err := mdb.ListCollectionSpecifications(ctx,
		bson.M{
			"name": "history",
		},
	)
	assert.NoError(t, err)
	if len(colls) != 0 {
		t.Log(colls[0].Type)
	}
}

func TestMgoCreateValidate(t *testing.T) {
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
	err := mdb.CreateCollection(ctx, "projects", option)
	assert.NoError(t, err)
	_, err = mdb.Collection("projects").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"namespace", 1}},
			Options: options.Index().SetName("idx_namespace").SetUnique(true),
		},
	})
	assert.NoError(t, err)
}

func TestMgoMockOrders(t *testing.T) {
	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(100, func(i interface{}) {
		_, err := mdb.Collection("orders").InsertMany(context.TODO(), i.([]interface{}))
		assert.NoError(t, err)
		wg.Done()
	})
	assert.NoError(t, err)
	defer p.Release()
	for n := 0; n < 100; n++ {
		wg.Add(1)
		orders := make([]interface{}, 10000)
		for i := 0; i < 10000; i++ {
			var data Order
			err = faker.FakeData(&data)
			assert.NoError(t, err)
			orders[i] = data
		}
		_ = p.Invoke(orders)
	}
	wg.Wait()
}

func TestMgoAvg(t *testing.T) {
	var avg []bson.M
	ctx := context.TODO()
	start := time.Now()
	c, err := mdb.Collection("orders").Aggregate(ctx, mongo.Pipeline{
		{
			{"$group", bson.D{
				{"_id", nil},
				{"avg", bson.D{{"$avg", "$price"}}},
			}},
		},
	})
	assert.NoError(t, err)
	err = c.All(ctx, &avg)
	assert.NoError(t, err)
	t.Log(time.Since(start))
	t.Log(avg)
}
