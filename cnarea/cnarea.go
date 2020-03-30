// Package cnarea 提供中国行政区域相关方法
//
// 数据库文件：../sql/cnarea20181031.sql.tar.gz
// 新建db：cnarea，导入数据
package cnarea

import (
	"fmt"

	// use mysql db
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Level 行政区域level字段
type Level int

const (
	// ProvinceLevel 省级、直辖市级、特别行政区级
	ProvinceLevel Level = iota
	// CityLevel 市级、直辖区级
	CityLevel
	// DistrictLevel 区级、县级
	DistrictLevel
	// TownLevel 乡镇级
	TownLevel
	// CommunityLevel 社区居委会级
	CommunityLevel
)

// NewDB 创建sqlx DB实例
func NewDB(host string, port int, user, passwd string) (*sqlx.DB, error) {
	mydsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/cnarea?charset=utf8&parseTime=true", user, passwd, host, port)
	db, err := sqlx.Open("mysql", mydsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Province 查询指定省级、直辖市级区域、特别行政区级信息
func Province(db *sqlx.DB, shortName string) (Area, error) {
	area := Area{}
	if err := db.Get(&area, "select * from cnarea_2018 where level = ? and short_name = ?", ProvinceLevel, shortName); err != nil {
		return area, err
	}
	return area, nil
}

// Provinces 查询指定省级、直辖市级区域、特别行政区级列表
func Provinces(db *sqlx.DB) ([]Area, error) {
	areas := []Area{}
	if err := db.Select(&areas, "select * from cnarea_2018 where level = ?", ProvinceLevel); err != nil {
		return nil, err
	}
	return areas, nil
}

// City 查询指定市级、直辖区级信息
func City(db *sqlx.DB, shortName string) (Area, error) {
	area := Area{}
	if err := db.Get(&area, "select * from cnarea_2018 where level = ? and short_name = ?", CityLevel, shortName); err != nil {
		return area, err
	}
	return area, nil
}

// Cities 查询指定省级、直辖市级区域、特别行政区级下的城市列表
func Cities(db *sqlx.DB, proviceShortName string) ([]Area, error) {
	areas := []Area{}
	provinceArea, err := Province(db, proviceShortName)
	if err != nil {
		return nil, err
	}
	if err := db.Select(&areas, "select * from cnarea_2018 where level = ? and parent_code = ?", CityLevel, provinceArea.AreaCode); err != nil {
		return nil, err
	}
	return areas, nil
}

// Districts 查询指定市级下的区列表
func Districts(db *sqlx.DB, cityShortName string) ([]Area, error) {
	areas := []Area{}
	cityArea, err := City(db, cityShortName)
	if err != nil {
		return nil, err
	}
	if err := db.Select(&areas, "select * from cnarea_2018 where level = ? and parent_code = ?", DistrictLevel, cityArea.AreaCode); err != nil {
		return nil, err
	}
	return areas, nil
}
