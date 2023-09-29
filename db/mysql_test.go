package db

import (
	"context"
	"database/sql"
	"github.com/go-faker/faker/v4"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMySQLCreate(t *testing.T) {
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
		Table("department").
		Model(&data).
		Exec(ctx)
	assert.NoError(t, err)

	t.Log(r.LastInsertId())
	t.Log(r.RowsAffected())
}

func TestMySQLFind(t *testing.T) {
	ctx := context.TODO()
	var data []map[string]interface{}

	if err := db.NewSelect().
		Table("department").
		Scan(ctx, &data); err != nil {
		t.Error(err)
	}

	t.Log(data)
}

func TestMySQLFindOne(t *testing.T) {
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
