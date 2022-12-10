package postgres

import (
	"bufio"
	"development/common"
	"development/postgres/model"
	"encoding/csv"
	"github.com/alexedwards/argon2id"
	"github.com/go-faker/faker/v4"
	"github.com/panjf2000/ants/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var values *common.Values
var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config/config.yml"); err != nil {
		log.Fatalln(err)
	}
	if db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  values.POSTGRES,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	if err := db.AutoMigrate(model.User{}); err != nil {
		t.Error(err)
	}
	hash, err := argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams)
	if err != nil {
		t.Error(err)
	}
	if err = db.Create(&model.User{
		Username: "weplanx",
		Password: hash,
		Email:    "zhangtqx@vip.qq.com",
	}).Error; err != nil {
		t.Error(err)
	}
}

func TestMockCity(t *testing.T) {
	var err error
	var f *os.File
	if f, err = os.Open("../assets/cities.csv"); err != nil {
		t.Error(err)
	}

	if err = db.AutoMigrate(&model.City{}); err != nil {
		t.Error(err)
	}

	r := csv.NewReader(f)
	first := true
	var cities []model.City
	for {
		var record []string
		if record, err = r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				t.Error(err)
			}
		}
		if first {
			first = false
			continue
		}
		latitude := float64(0)
		if record[8] != "" {
			if latitude, err = strconv.ParseFloat(record[8], 64); err != nil {
				t.Error(err)
			}
		}
		longitude := float64(0)
		if record[9] != "" {
			if longitude, err = strconv.ParseFloat(record[9], 64); err != nil {
				t.Error(err)
			}
		}
		cities = append(cities, model.City{
			Name:        record[1],
			CountryCode: record[6],
			StateCode:   record[3],
			Latitude:    latitude,
			Longitude:   longitude,
		})
	}
	if err = db.CreateInBatches(cities, 5000).Error; err != nil {
		t.Error(err)
	}
}

func TestMockKVCity(t *testing.T) {
	var err error
	var f *os.File
	if f, err = os.Open("../assets/cities.csv"); err != nil {
		t.Error(err)
	}

	if err = db.AutoMigrate(&model.KVCity{}); err != nil {
		t.Error(err)
	}

	r := csv.NewReader(f)
	first := true
	var cities []model.KVCity
	for {
		var record []string
		if record, err = r.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				t.Error(err)
			}
		}
		if first {
			first = false
			continue
		}
		latitude := float64(0)
		if record[8] != "" {
			if latitude, err = strconv.ParseFloat(record[8], 64); err != nil {
				t.Error(err)
			}
		}
		longitude := float64(0)
		if record[9] != "" {
			if longitude, err = strconv.ParseFloat(record[9], 64); err != nil {
				t.Error(err)
			}
		}
		cities = append(cities, model.KVCity{
			Value: model.CityValue{
				Name:        record[1],
				CountryCode: record[6],
				StateCode:   record[3],
				Latitude:    latitude,
				Longitude:   longitude,
			},
		})
	}
	if err = db.CreateInBatches(cities, 5000).Error; err != nil {
		t.Error(err)
	}
}

func TestQueryCity(t *testing.T) {
	var data model.City
	if err := db.Debug().
		Where("name = ?", "Pogradec").
		Take(&data).Error; err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestQueryKVCity(t *testing.T) {
	var data model.KVCity
	if err := db.Debug().
		Where("value ->> 'name' = ?", "Pogradec").
		Take(&data).Error; err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func isZero(value string) string {
	if value == "0" {
		return ""
	}
	return value
}

func ip2Dec(value string) uint64 {
	ip := uint64(0)
	for k, v := range strings.Split(value, ".") {
		n, _ := strconv.ParseUint(v, 10, 64)
		ip |= n << ((3 - uint64(k)) << 3)
	}
	return ip
}

func TestMockIp(t *testing.T) {
	var err error
	var f *os.File
	if f, err = os.Open("../assets/ip.merge.txt"); err != nil {
		t.Error(err)
	}

	if err = db.AutoMigrate(&model.Ipv4{}); err != nil {
		t.Error(err)
	}

	var bulk []model.Ipv4
	r := bufio.NewReader(f)
	for {
		var s string
		if s, err = r.ReadString('\n'); err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				t.Error(err)
			}
		}
		row := strings.TrimSpace(s)
		if row == "" {
			continue
		}
		v := strings.Split(row, "|")
		bulk = append(bulk, model.Ipv4{
			Start:    ip2Dec(v[0]),
			End:      ip2Dec(v[1]),
			Country:  isZero(v[2]),
			Province: isZero(v[4]),
			City:     isZero(v[5]),
			ISP:      isZero(v[6]),
		})
	}

	if err = db.CreateInBatches(bulk, 5000).Error; err != nil {
		t.Error(err)
	}
}

func TestMockKVIp(t *testing.T) {
	var err error
	var f *os.File
	if f, err = os.Open("../assets/ip.merge.txt"); err != nil {
		t.Error(err)
	}

	if err = db.AutoMigrate(&model.KVIpv4{}); err != nil {
		t.Error(err)
	}

	var bulk []model.KVIpv4
	r := bufio.NewReader(f)
	for {
		var s string
		if s, err = r.ReadString('\n'); err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				t.Error(err)
			}
		}
		row := strings.TrimSpace(s)
		if row == "" {
			continue
		}
		v := strings.Split(row, "|")
		bulk = append(bulk, model.KVIpv4{
			Value: model.IpValue{
				Start:    ip2Dec(v[0]),
				End:      ip2Dec(v[1]),
				Country:  isZero(v[2]),
				Province: isZero(v[4]),
				City:     isZero(v[5]),
				ISP:      isZero(v[6]),
			},
		})
	}

	if err = db.CreateInBatches(bulk, 5000).Error; err != nil {
		t.Error(err)
	}
}

func TestMockOrder(t *testing.T) {
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	p, err := ants.NewPoolWithFunc(100, func(i interface{}) {
		if err := db.CreateInBatches(i.([]model.Order), 2000).Error; err != nil {
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
		orders := make([]model.Order, 10000)
		for i := 0; i < 10000; i++ {
			if err := faker.FakeData(&orders[i]); err != nil {
				t.Error(err)
			}
		}
		_ = p.Invoke(orders)
	}
	wg.Wait()
}
