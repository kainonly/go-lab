package mysql

import (
	"context"
	"database/sql"
	"github.com/go-faker/faker/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"go-lab/common"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var values *common.Values
var db *bun.DB

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

func TestMySQLCreate(t *testing.T) {
	ctx := context.TODO()
	data := map[string]interface{}{
		"name":        "测试",
		"description": "部门",
		"schema": []M{
			{"key": "asd"},
		},
		"create_time": time.Now(),
		"update_time": time.Now(),
	}

	r, err := db.NewInsert().
		Table("department").
		Model(&data).
		Exec(ctx)
	assert.NoError(t, err)

	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}

func TestMySQLFind(t *testing.T) {
	ctx := context.TODO()
	var data []M

	err := db.NewSelect().
		Table("department").
		Scan(ctx, &data)
	assert.NoError(t, err)

	t.Log(data)
}

func TestMySQLFindOne(t *testing.T) {
	ctx := context.TODO()
	var data map[string]interface{}

	err := db.NewSelect().
		Table("departments").
		Where(`id = 1`).
		Scan(ctx, &data)
	assert.NoError(t, err)

	t.Log(data)
}

func TestMySQLUpdate(t *testing.T) {
	ctx := context.TODO()
	data := map[string]interface{}{
		"name":        "测试123",
		"update_time": time.Now(),
	}

	r, err := db.NewUpdate().
		Table("department").
		Where(`id = 1`).
		Model(&data).
		Exec(ctx)
	assert.NoError(t, err)

	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}

func TestMySQLDelete(t *testing.T) {
	ctx := context.TODO()

	r, err := db.NewDelete().
		Table("departments").
		Where(`id = 2`).
		Exec(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}

func TestMySQLTransaction(t *testing.T) {
	ctx := context.TODO()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})

	assert.NoError(t, err)
	r, err := tx.NewDelete().
		Table("department").
		Where(`id = 7`).
		Exec(ctx)
	assert.NoError(t, err)

	tx.Rollback()
	t.Log(r)
}

func TestMySQLMock(t *testing.T) {
	ctx := context.TODO()
	err := db.ResetModel(ctx, (*IOrder)(nil))
	assert.NoError(t, err)
	var wg sync.WaitGroup
	var p *ants.PoolWithFunc
	p, err = ants.NewPoolWithFunc(1000, func(i interface{}) {
		_, err = db.NewInsert().Model(i.(*[]IOrder)).Exec(ctx)
		assert.NoError(t, err)
		wg.Done()
	})
	assert.NoError(t, err)
	defer p.Release()

	for w := 0; w < 100; w++ {
		wg.Add(1)
		orders := make([]IOrder, 10000)
		for i := 0; i < 10000; i++ {
			var data IOrder
			err = faker.FakeData(&data)
			assert.NoError(t, err)
			orders[i] = data
		}
		_ = p.Invoke(&orders)
	}

	wg.Wait()
}

type IProject struct {
	bun.BaseModel `bun:"table:project"`
	ID            uint64    `bun:"id,pk,autoincrement" faker:"-"`
	Name          string    `bun:"type:varchar(50),notnull"`
	Namespace     string    `bun:"type:varchar(20),notnull,unique"`
	Meta          M         `bun:"type:json,notnull"`
	CreateTime    time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdateTime    time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func TestMySQLRestInit(t *testing.T) {
	ctx := context.TODO()
	err := db.ResetModel(ctx, (*IProject)(nil))
	assert.NoError(t, err)
}

func TestMySQLRestCreate(t *testing.T) {
	ctx := context.TODO()

	data := M{
		"name":      "测试",
		"namespace": "beta",
		"meta": M{
			"n": 123,
			"p": "xxx",
		},
		"create_time": time.Now(),
		"update_time": time.Now(),
	}

	r, err := db.NewInsert().
		Table("project").
		Model(&data).
		Exec(ctx)

	assert.NoError(t, err)

	t.Log(r)
}

type IPicture struct {
	bun.BaseModel `bun:"table:picture"`
	ID            int    `bun:"id,pk,autoincrement" faker:"-"`
	Tid           []int  `bun:"type:json" faker:"-"`
	Name          string `bun:"type:varchar(50)" faker:"name"`
}

func RandomTid() []int {
	n := rand.Intn(10)
	return rand.Perm(10)[:n]
}

func TestMySQLMockJson(t *testing.T) {
	ctx := context.TODO()
	err := db.ResetModel(ctx, (*IPicture)(nil))
	assert.NoError(t, err)
	var wg sync.WaitGroup
	var p *ants.PoolWithFunc
	p, err = ants.NewPoolWithFunc(1000, func(i interface{}) {
		_, err = db.NewInsert().Model(i.(*[]IPicture)).Exec(ctx)
		assert.NoError(t, err)
		wg.Done()
	})
	assert.NoError(t, err)
	defer p.Release()

	for w := 0; w < 1; w++ {
		wg.Add(1)
		pictures := make([]IPicture, 10000)
		for i := 0; i < 10000; i++ {
			var data IPicture
			err = faker.FakeData(&data)
			data.Tid = RandomTid()
			assert.NoError(t, err)
			pictures[i] = data
		}
		_ = p.Invoke(&pictures)
	}

	wg.Wait()
}
