package weiboclock

import (
	"fmt"
	"io"
	"log"

	// 导入statik生成的代码
	_ "cuitclock/statik"

	"github.com/axiaoxin-com/cronweibo"
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
	// 今日使用进度
	dayProcessBar := DayProgressBar(hour)
	// 天文信息
	cityAstroInfo, err := CityAstroInfo(viper.GetString("wttrin.location"), now)
	if err != nil {
		log.Println("[ERROR] weiboclock tollJob CityAstroInfo error", err)
	}

	text := fmt.Sprintf("%s %d点啦%s %s\n\n"+
		"今日进度:\n%s\n\n"+
		"%s%s\n\n",
		ClockEmoji[oclock], oclock, TollTail(oclock), emotion,
		dayProcessBar,
		WttrInLine, cityAstroInfo)

	// 生成图片内容
	pic, err := PicReader(viper.GetString("weiboclock.pic_path"), now)
	if err != nil {
		log.Println("[ERROR] weiboclock toll error:", err)
		// 有error也不影响发送，获取图片失败就不发图片
	}
	return text, pic
}
