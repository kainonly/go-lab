package postgres

import (
	"context"
	"database/sql"
	"development/common"
	"github.com/go-faker/faker/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

var values *common.Values
var db *bun.DB

func TestMain(m *testing.M) {
	var err error
	os.Chdir("../../")
	if values, err = common.LoadValues("./config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	sqldb, err := sql.Open("mysql", values.MYSQL)
	if err != nil {
		panic(err)
	}
	db = bun.NewDB(sqldb, mysqldialect.New())
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	data := map[string]interface{}{
		"name":        "测试",
		"description": "部门",
		"schema": []map[string]interface{}{
			{"key": "asd"},
		},
		"create_time": time.Now(),
		"update_time": time.Now(),
	}

	r, err := db.NewInsert().
		Table("departments").
		Model(&data).
		Exec(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}

func TestFind(t *testing.T) {
	ctx := context.TODO()
	var data []map[string]interface{}

	if err := db.NewSelect().
		Table("departments").
		Scan(ctx, &data); err != nil {
		t.Error(err)
	}

	t.Log(data)
}

func TestFindOne(t *testing.T) {
	ctx := context.TODO()
	var data map[string]interface{}

	if err := db.NewSelect().
		Table("departments").
		Where(`id = 1`).
		Scan(ctx, &data); err != nil {
		t.Error(err)
	}

	t.Log(data)
}

func TestUpdate(t *testing.T) {
	ctx := context.TODO()
	data := map[string]interface{}{
		"name":        "测试123",
		"update_time": time.Now(),
	}

	r, err := db.NewUpdate().
		Table("departments").
		Where(`id = 1`).
		Model(&data).
		Exec(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}

func TestDelete(t *testing.T) {
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

func TestTransaction(t *testing.T) {
	ctx := context.TODO()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})

	assert.NoError(t, err)
	r, err := tx.NewDelete().
		Table("departments").
		Where(`id = 7`).
		Exec(ctx)
	if err != nil {
		t.Error(err)
	}

	tx.Rollback()
	t.Log(r)
}

type Order struct {
	ID          uint64  `bun:"id,pk,autoincrement" faker:"-"`
	No          string  `bun:"type:varchar(255)" faker:"cc_number"`
	Name        string  `bun:"type:varchar(255)" faker:"name"`
	Description string  `bun:"type:text" faker:"paragraph"`
	Account     string  `bun:"type:varchar(255)" faker:"username"`
	Customer    string  `bun:"type:varchar(255)" faker:"name"`
	Email       string  `bun:"type:varchar(255)" faker:"email"`
	Phone       string  `bun:"type:varchar(255)" faker:"phone_number"`
	Address     string  `bun:"type:varchar(255)" faker:"sentence"`
	Price       float64 `bun:"type:decimal" faker:"amount"`
}

func TestMockOrder(t *testing.T) {
	ctx := context.TODO()
	err := db.ResetModel(ctx, (*Order)(nil))
	var wg sync.WaitGroup
	var p *ants.PoolWithFunc
	p, err = ants.NewPoolWithFunc(1000, func(i interface{}) {
		_, err = db.NewInsert().Model(i.(*[]Order)).Exec(ctx)
		assert.NoError(t, err)
		wg.Done()
	})
	assert.NoError(t, err)
	defer p.Release()

	for w := 0; w < 150; w++ {
		wg.Add(1)
		orders := make([]Order, 10000)
		for i := 0; i < 10000; i++ {
			err = faker.FakeData(&orders[i])
			assert.NoError(t, err)
		}
		_ = p.Invoke(&orders)
	}

	wg.Wait()
}

type OrderY Order

func TestMockOrderY(t *testing.T) {
	ctx := context.TODO()
	err := db.ResetModel(ctx, (*OrderY)(nil))
	var wg sync.WaitGroup
	var p *ants.PoolWithFunc
	p, err = ants.NewPoolWithFunc(1000, func(i interface{}) {
		_, err = db.NewInsert().Model(i.(*[]OrderY)).Exec(ctx)
		assert.NoError(t, err)
		wg.Done()
	})
	assert.NoError(t, err)
	defer p.Release()

	for w := 0; w < 150*12; w++ {
		wg.Add(1)
		orders := make([]OrderY, 10000)
		for i := 0; i < 10000; i++ {
			err = faker.FakeData(&orders[i])
			assert.NoError(t, err)
		}
		_ = p.Invoke(&orders)
	}

	wg.Wait()
}
