package cnarea

import (
	"log"
	"testing"

	"github.com/spf13/viper"
)

func new() *Query {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	q, err := NewQuery(viper.GetString("mysql.host"), viper.GetInt("mysql.port"), viper.GetString("mysql.user"), viper.GetString("mysql.passwd"))
	if err != nil {
		log.Fatal(err)
	}
	return q
}

func TestProvinceLevelArea(t *testing.T) {
	q := new()
	defer q.Close()
	// 查询四川省区域信息
	p, err := q.ProvinceLevelArea("四川")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestProvinceLevelAreas(t *testing.T) {
	q := new()
	defer q.Close()
	// 查询全国所有省+直辖市+特别行政区区域信息列表
	p, err := q.ProvinceLevelAreas()
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestCityLevelArea(t *testing.T) {
	q := new()
	defer q.Close()
	// 查询成都市区域信息
	p, err := q.CityLevelArea("成都")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestCityLevelAreas(t *testing.T) {
	q := new()
	defer q.Close()
	// 查询四川省所有市级区域信息列表
	p, err := q.CityLevelAreas("四川")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func TestDistrictLevelAreas(t *testing.T) {
	q := new()
	defer q.Close()
	// 查询成都市所有区县级区域信息列表
	p, err := q.DistrictLevelAreas("成都")
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}
