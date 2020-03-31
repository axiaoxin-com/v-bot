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

// Query cnarea数据查询实例
type Query struct {
	db *sqlx.DB
}

// Close 关闭db连接
func (q *Query) Close() {
	q.db.Close()
}

// DB 获取sqlx db 实例
func (q *Query) DB() *sqlx.DB {
	return q.db
}

// NewQuery 创建sqlx DB实例
func NewQuery(host string, port int, user, passwd string) (*Query, error) {
	mydsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/cnarea?charset=utf8&parseTime=true", user, passwd, host, port)
	db, err := sqlx.Open("mysql", mydsn)
	if err != nil {
		return nil, err
	}
	return &Query{db: db}, nil
}

// ProvinceLevelArea 查询指定省级、直辖市级区域、特别行政区级的区域信息
func (q *Query) ProvinceLevelArea(areaShortName string) (Area, error) {
	area := Area{}
	if err := q.db.Get(&area, "select * from cnarea_2018 where level = ? and short_name = ?", ProvinceLevel, areaShortName); err != nil {
		return area, err
	}
	return area, nil
}

// ProvinceLevelAreas 查询所有省级+直辖市级+特别行政区级的区域信息列表
func (q *Query) ProvinceLevelAreas() ([]Area, error) {
	areas := []Area{}
	if err := q.db.Select(&areas, "select * from cnarea_2018 where level = ?", ProvinceLevel); err != nil {
		return nil, err
	}
	return areas, nil
}

// CityLevelArea 查询指定市级、直辖区级的区域信息
func (q *Query) CityLevelArea(areaShortName string) (Area, error) {
	area := Area{}
	if err := q.db.Get(&area, "select * from cnarea_2018 where level = ? and short_name = ?", CityLevel, areaShortName); err != nil {
		return area, err
	}
	return area, nil
}

// CityLevelAreas 查询指定省级、直辖市级、特别行政区级区域下的所有市级、直辖区级的区域信息列表
func (q *Query) CityLevelAreas(areaShortName string) ([]Area, error) {
	areas := []Area{}
	if err := q.db.Select(&areas, "select * from cnarea_2018 where level = ? and parent_code = (select area_code from cnarea_2018 where level = ? and short_name = ?)", CityLevel, ProvinceLevel, areaShortName); err != nil {
		return nil, err
	}
	return areas, nil
}

// DistrictLevelAreas 查询指定市级区域下的所有区级区域信息列表
func (q *Query) DistrictLevelAreas(areaShortName string) ([]Area, error) {
	areas := []Area{}
	if err := q.db.Select(&areas, "select * from cnarea_2018 where level = ? and parent_code = (select area_code from cnarea_2018 where level = ? and short_name = ?)", DistrictLevel, CityLevel, areaShortName); err != nil {
		return nil, err
	}
	return areas, nil
}
