package weiboclock

import (
	"cuitclock/cnarea"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/starainrt/astro"
)

// Lunar 查询今日农历日期 📅
func Lunar(t time.Time) string {
	// 农历
	_, _, _, lunar := astro.Lunar(t.Year(), int(t.Month()), t.Day())
	return lunar
}

// Sunrise 指定经纬度的今日日出时间 🌅
func Sunrise(lng, lat float64, t time.Time) time.Time {
	_, offset := t.Zone()
	zone := offset / 60 / 60
	// 日出
	sunrise, _ := astro.SunRiseTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	return sunrise
}

// Sunset 指定经纬度的今日日落时间 🌄
func Sunset(lng, lat float64, t time.Time) time.Time {
	_, offset := t.Zone()
	zone := offset / 60 / 60
	sunset, _ := astro.SunDownTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	return sunset
}

// Moonrise 指定经纬度的今日月出时间
func Moonrise(lng, lat float64, t time.Time) time.Time {
	_, offset := t.Zone()
	zone := offset / 60 / 60
	moonrise, _ := astro.MoonRiseTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	return moonrise
}

// Moonset 指定经纬度的今日月落时间
func Moonset(lng, lat float64, t time.Time) time.Time {
	_, offset := t.Zone()
	zone := offset / 60 / 60
	moonset, _ := astro.MoonDownTime(astro.Date2JDE(t), lng, lat, float64(zone), true)
	return moonset
}

// CityAstroInfo 根据城市名称获取当地指定时间天文信息
func CityAstroInfo(cityname string, t time.Time) (string, error) {
	q, err := cnarea.NewQuery(viper.GetString("mysql.host"), viper.GetInt("mysql.port"), viper.GetString("mysql.user"), viper.GetString("mysql.passwd"))
	if err != nil {
		return "", err
	}
	city, err := q.CityLevelArea(cityname)
	if err != nil {
		return "", err
	}
	info := fmt.Sprintf("农历📅 %s\n"+
		"日出🌅 %s\n"+
		"日落🌄 %s",
		Lunar(t),
		Sunrise(city.Lng, city.Lat, t).Format("15:04:05"),
		Sunset(city.Lng, city.Lat, t).Format("15:04:05"),
	)
	return info, nil
}
