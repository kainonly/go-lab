package postgres

import (
	"context"
	"development/common"
	"github.com/go-faker/faker/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

var values *common.Values
var db *gorm.DB

type Department struct {
	Id          int64     `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar"`
	Description string    `gorm:"type:varchar"`
	CreateTime  time.Time `gorm:"autoCreateTime"`
	UpdateTime  time.Time `gorm:"autoUpdateTime"`
}

func TestMain(m *testing.M) {
	var err error
	os.Chdir("../")
	if values, err = common.LoadValues("./config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	if db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  values.POSTGRES,
		PreferSimpleProtocol: true,
	}), &gorm.Config{}); err != nil {
		log.Fatalln(err)
	}
	if err = db.AutoMigrate(&Department{}); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	data := Department{
		Name:        faker.Name(),
		Description: faker.Paragraph(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err := db.WithContext(ctx).Create(&data).Error
	assert.NoError(t, err)
	t.Log(data.Id)
}

//	func TestFind(t *testing.T) {
//		ctx := context.TODO()
//		var data []map[string]interface{}
//
//		if err := db.NewSelect().
//			Table("departments").
//			Scan(ctx, &data); err != nil {
//			t.Error(err)
//		}
//
//		t.Log(data)
//	}
//
//	func TestFindOne(t *testing.T) {
//		ctx := context.TODO()
//		var data map[string]interface{}
//
//		if err := db.NewSelect().
//			Table("departments").
//			Where(`id = 1`).
//			Scan(ctx, &data); err != nil {
//			t.Error(err)
//		}
//
//		t.Log(data)
//	}
//
//	func TestUpdate(t *testing.T) {
//		ctx := context.TODO()
//		data := map[string]interface{}{
//			"name":        "测试123",
//			"update_time": time.Now(),
//		}
//
//		r, err := db.NewUpdate().
//			Table("departments").
//			Where(`id = 1`).
//			Model(&data).
//			Exec(ctx)
//		if err != nil {
//			t.Error(err)
//		}
//
//		t.Log(r.LastInsertId())
//		t.Log(r.RowsAffected())
//	}
//
//	func TestDelete(t *testing.T) {
//		ctx := context.TODO()
//
//		r, err := db.NewDelete().
//			Table("departments").
//			Where(`id = 2`).
//			Exec(ctx)
//		if err != nil {
//			t.Error(err)
//		}
//
//		t.Log(r.LastInsertId())
//		t.Log(r.RowsAffected())
//	}
//
//	func TestTx(t *testing.T) {
//		ctx := context.TODO()
//
//		tx, err := db.BeginTx(ctx, &sql.TxOptions{})
//
//		assert.NoError(t, err)
//		r, err := tx.NewDelete().
//			Table("departments").
//			Where(`id = 7`).
//			Exec(ctx)
//		if err != nil {
//			t.Error(err)
//		}
//
//		tx.Rollback()
//		t.Log(r)
//	}
//
//	type Order struct {
//		ID          uint64  `bun:"id,pk,autoincrement" faker:"-"`
//		No          string  `bun:"type:varchar" faker:"cc_number"`
//		Name        string  `bun:"type:varchar" faker:"name"`
//		Description string  `bun:"type:text" faker:"paragraph"`
//		Account     string  `bun:"type:varchar" faker:"username"`
//		Customer    string  `bun:"type:varchar" faker:"name"`
//		Email       string  `bun:"type:varchar" faker:"email"`
//		Phone       string  `bun:"type:varchar" faker:"phone_number"`
//		Address     string  `bun:"type:varchar" faker:"sentence"`
//		Price       float64 `bun:"type:decimal" faker:"amount"`
//		CreateTime  time.Time
//		UpdateTime  time.Time
//	}
//
//	func TestMockOrder(t *testing.T) {
//		ctx := context.TODO()
//		err := db.ResetModel(ctx, (*Order)(nil))
//		var wg sync.WaitGroup
//		var p *ants.PoolWithFunc
//		p, err = ants.NewPoolWithFunc(1000, func(i interface{}) {
//			_, err = db.NewInsert().Model(i.(*[]Order)).Exec(ctx)
//			assert.NoError(t, err)
//			wg.Done()
//		})
//		assert.NoError(t, err)
//		defer p.Release()
//
//		for w := 0; w < 100; w++ {
//			wg.Add(1)
//			orders := make([]Order, 10000)
//			for i := 0; i < 10000; i++ {
//				var data Order
//				err = faker.FakeData(&data)
//				assert.NoError(t, err)
//				data.CreateTime, _ = time.Parse(`2006-01-02 15:04:05`, faker.Timestamp())
//				data.UpdateTime = data.CreateTime.Add(time.Hour * 24)
//				orders[i] = data
//			}
//			_ = p.Invoke(&orders)
//		}
//
//		wg.Wait()
//	}

type Consumption struct {
	Customer string
	Total    float64
}

func TestConsumption(t *testing.T) {
	ctx := context.TODO()
	var data []Consumption
	err := db.WithContext(ctx).
		Table("orders").
		Select("customer,sum(price) as total").
		Group("customer").
		Order("total desc").
		Limit(10).
		Find(&data).Error

	assert.NoError(t, err)
	t.Log(data)
}
