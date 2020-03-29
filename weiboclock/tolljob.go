package weiboclock

import (
	"fmt"
	"io"
	"log"

	// 导入statik生成的代码
	_ "cuitclock/statik"

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
	hour := now.Hour()
	oclock := hour
	// 12小时制处理
	if hour > 12 {
		oclock = hour - 12
	} else if hour == 0 {
		oclock = 12
	}
	weather, err := wttrin.Line(viper.GetString("wttrin.lang"), viper.GetString("wttrin.location"), "")
	if err != nil {
		log.Println("[ERROR] tollRun get weather error", err)
	}
	text := fmt.Sprintf("%d点啦~ %s %s\n\n%s\n", oclock, emotion, TollVoice(oclock), weather)

	// 生成图片内容
	pic, err := clock.PicReader(viper.GetString("weiboclock.pic_path"), hour)
	if err != nil {
		log.Println("[ERROR] weiboclock toll error:", err)
		// 有error也不影响发送，获取图片失败就不发图片
	}
	return text, pic
}
