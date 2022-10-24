package mongodb

import (
	"context"
	"development/common"
	"development/mongodb/model"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"os"
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

	docs := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		var doc model.Record
		if err := faker.FakeData(&doc); err != nil {
			t.Error(err)
		}
		docs[i] = doc
	}

	if _, err := db.Collection("history").InsertMany(ctx, docs); err != nil {
		t.Error(err)
	}
}
