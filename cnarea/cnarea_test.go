package cnarea

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func new() *sqlx.DB {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	db, err := NewDB(viper.GetString("mysql.host"), viper.GetInt("mysql.port"), viper.GetString("mysql.user"), viper.GetString("mysql.passwd"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestProvince(t *testing.T) {
	db := new()
	defer db.Close()
	p, err := Province(db, "四川")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestProvinces(t *testing.T) {
	db := new()
	defer db.Close()
	p, err := Provinces(db)
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestCity(t *testing.T) {
	db := new()
	defer db.Close()
	p, err := City(db, "成都")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestCities(t *testing.T) {
	db := new()
	defer db.Close()
	p, err := Cities(db, "四川")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestDistricts(t *testing.T) {
	db := new()
	defer db.Close()
	p, err := Districts(db, "成都")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}
