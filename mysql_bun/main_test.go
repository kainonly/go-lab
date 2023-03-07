package postgres

import (
	"context"
	"database/sql"
	"development/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var db *bun.DB

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
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
