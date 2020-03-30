package cnarea

// Area 中国行政地区表
type Area struct {
	ID int `json:"id" db:"id"`
	//Level 层级
	Level int8 `json:"level" db:"level"`
	//ParentCode 父级行政代码
	ParentCode int64 `json:"parent_code" db:"parent_code"`
	//AreaCode 行政代码
	AreaCode int64 `json:"area_code" db:"area_code"`
	//ZIPCode 邮政编码
	ZIPCode int `json:"zip_code" db:"zip_code"`
	//CityCode 区号
	CityCode string `json:"city_code" db:"city_code"`
	//Name 名称
	Name string `json:"name" db:"name"`
	//ShortName 简称
	ShortName string `json:"short_name" db:"short_name"`
	//MergerName 组合名
	MergerName string `json:"merger_name" db:"merger_name"`
	//Pinyin 拼音
	Pinyin string `json:"pinyin" db:"pinyin"`
	//Lng 经度
	Lng float64 `json:"lng" db:"lng"`
	//Lat 纬度
	Lat float64 `json:"lat" db:"lat"`
}

// TableName cnarea_2018
func (t *Area) TableName() string {
	return "cnarea_2018"
}
