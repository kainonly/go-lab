package mysql

import (
	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"context"
	"database/sql"
	"development/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var values *common.Values
var db *sql.DB
var drv migrate.Driver
var sch *schema.Schema
var r *schema.Realm

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	if db, err = sql.Open("mysql", values.MYSQL); err != nil {
		log.Fatalln(err)
	}
	if drv, err = mysql.Open(db); err != nil {
		log.Fatalln(err)
	}
	if sch, err = drv.InspectSchema(context.TODO(), "example", &schema.InspectOptions{
		Tables: []string{"roles"},
	}); err != nil {
		log.Fatalln(err)
	}
	if r, err = drv.InspectRealm(context.TODO(), &schema.InspectRealmOption{
		Schemas: []string{"example"},
	}); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestCreateTable(t *testing.T) {
	ctx := context.Background()
	to := schema.NewTable("roles").
		SetPrimaryKey(schema.NewPrimaryKey(
			schema.NewColumn("id").
				SetType(&schema.IntegerType{T: "int"}).
				SetNull(false)),
		).
		AddColumns(schema.NewColumn("name").
			SetType(&schema.StringType{T: "varchar", Size: 20}).
			SetNull(false),
		)
	example, ok := r.Schema("example")
	assert.True(t, ok)
	from, ok := example.Table("roles")
	assert.True(t, ok)
	changes, err := drv.TableDiff(from, to)
	assert.NoError(t, err)
	err = drv.ApplyChanges(ctx, changes)
	assert.NoError(t, err)
}
