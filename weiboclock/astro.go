package weiboclock

import (
	"cuitclock/cnarea"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/starainrt/astro"
)

// AstroInfo 天文信息
type AstroInfo struct {
	Lunar    string    // 今日农历日期
	Sunrise  time.Time // 指定经纬度的今日日出时间
	Sunset   time.Time // 指定经纬度的今日日落时间
	Moonrise time.Time // 指定经纬度的今日月出时间
	Moonset  time.Time // 指定经纬度的今日月落时间
}

// GetAstroInfo 获取指定经纬度当前的天文信息
func GetAstroInfo(lng, lat float64, t time.Time) *AstroInfo {
	_, offset := t.Zone()
	zone := offset / 60 / 60
	// 农历
	_, _, _, lunar := astro.Lunar(t.Year(), int(t.Month()), t.Day())
	// 日出
	sunrise, _ := astro.SunRiseTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	sunset, _ := astro.SunDownTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	moonrise, _ := astro.MoonRiseTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	moonset, _ := astro.MoonDownTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	return &AstroInfo{
		Lunar:    lunar,
		Sunrise:  sunrise,
		Sunset:   sunset,
		Moonrise: moonrise,
		Moonset:  moonset,
	}
}

func (a *AstroInfo) String() string {
	return fmt.Sprintf("农历📆 %s\n"+
		"日出🌅 %s\n"+
		"日落🌄 %s\n",
		//"月出🌃 %s\n"+
		//"月落🏙 %s",
		a.Lunar,
		a.Sunrise.Format("15:04:05"),
		a.Sunset.Format("15:04:05"),
		//a.Moonrise.Format("15:04:05"),
		//a.Moonset.Format("15:04:05"),
	)
}

// CityAstroInfo 根据城市名称获取当地指定时间天文信息
func CityAstroInfo(cityname string, t time.Time) (*AstroInfo, error) {
	q, err := cnarea.NewQuery(viper.GetString("mysql.host"), viper.GetInt("mysql.port"), viper.GetString("mysql.user"), viper.GetString("mysql.passwd"))
	if err != nil {
		return nil, err
	}
	city, err := q.CityLevelArea(cityname)
	if err != nil {
		return nil, err
	}
	info := GetAstroInfo(city.Lng, city.Lat, t)
	return info, nil
}
