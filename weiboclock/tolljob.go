package weiboclock

import (
	"fmt"
	"io"
	"log"
	"unicode/utf8"

	"v-bot/progressbar"

	// 导入statik生成的代码
	_ "v-bot/statik"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/axiaoxin-com/wttrin"
	"github.com/spf13/viper"
)

// 返回整点报时任务
func (clock *WeiboClock) tollJob() cronweibo.WeiboJob {
	return cronweibo.WeiboJob{
		Name:     "toll",
		Schedule: "@hourly",
		Run:      clock.tollRun,
	}
}

// 返回整点报时的文字和图片，用于创建job
func (clock *WeiboClock) tollRun() (string, io.Reader) {
	// 生成文本内容
	now := clock.cronWeibo.Now()
	emotion := PickOneEmotion()
	log.Println("[DEBUG] tollRun picked emotion", emotion)
	// 24 小时制时刻
	hour := now.Hour()
	// 12 小时制时刻
	oclock := hour % 12
	if oclock == 0 {
		oclock = 12
	}
	// 今日使用进度
	dayProcessBar := progressbar.DayProgressBar(now)
	// 天文信息
	cityAstroInfo, err := CityAstroInfo(viper.GetString("weiboclock.wttrin_location"), now)
	if err != nil {
		log.Println("[ERROR] weiboclock tollJob CityAstroInfo error", err)
	}
	lang := viper.GetString("weiboclock.wttrin_lang")
	loc := viper.GetString("weiboclock.wttrin_location")
	format := "当前%l:\n天气%c %C\n温度🌡️ %t\n风速🌬️ %w\n湿度💦 %h\n月相🌑 +%M%m"
	weather, err := wttrin.Line(lang, loc, format)
	if err != nil {
		log.Println("[ERROR] weiboclock tollJob weather error", err)
	}

	text := fmt.Sprintf("%s %d点啦%s %s\n\n"+
		"你的今日使用进度:\n%s\n\n"+
		"%s"+
		"%s",
		ClockEmoji[oclock], oclock, TollTail(1), emotion,
		dayProcessBar,
		weather,
		cityAstroInfo,
	)
	log.Printf("[DEBUG] text:%s runecount:%d", text, utf8.RuneCountInString(text))
	// 生成图片内容
	pic, err := PicReader(viper.GetString("weiboclock.pic_path"), now)
	if err != nil {
		log.Println("[ERROR] weiboclock toll error:", err)
		// 有error也不影响发送，获取图片失败就不发图片
	}
	return text, pic
}
