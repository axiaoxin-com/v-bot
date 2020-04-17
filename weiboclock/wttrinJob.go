package weiboclock

import (
	"io"
	"log"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/axiaoxin-com/wttrin"
	"github.com/spf13/viper"
)

var (
	// WttrInLine 保存提前加载的天气文字信息
	WttrInLine string
	// WttrInImage 保存提前加载的天气图片内容
	WttrInImage io.ReadCloser
)

// 定时更新天气全局变量
func (clock *WeiboClock) wttrinJob() cronweibo.CronJob {
	return cronweibo.CronJob{
		Name:     "wttrin",
		Schedule: viper.GetString("weiboclock.wttrin_refresh_schedule"),
		Run:      clock.wttrinRun,
	}
}

// 提前加载天气信息
func (clock *WeiboClock) wttrinRun() {
	// reset
	WttrInLine = ""
	WttrInImage = nil

	// 默认在整点前 5 分钟更新天气
	viper.SetDefault("weiboclock.wttrin_refresh_schedule", "55 * * * *")
	lang := viper.GetString("weiboclock.wttrin_lang")
	loc := viper.GetString("weiboclock.wttrin_location")

	// 获取天气文本
	log.Println("[DEBUG] wttrinRun start getting Line weather")
	format := "当前%l:\n天气%c %C\n温度🌡️ %t\n风速🌬️ %w\n湿度💦 %h\n月相🌑 +%M%m"
	weather, err := wttrin.Line(lang, loc, format)
	if err == nil {
		WttrInLine = weather
		log.Println("[DEBUG] wttrinRun got the wttrin Line weather")
	} else {
		log.Println("[ERROR] wttrinRun get line weather error", err)
	}

	// 获取天气图片
	log.Println("[DEBUG] wttrinRun start getting Image weather")
	img, err := wttrin.Image(lang, loc)
	if err == nil {
		WttrInImage = img
		log.Println("[DEBUG] wttrinRun got the wttrin Image weather")
	} else {
		log.Println("[ERROR] wrttinRun get image weather error", err)
	}
}
