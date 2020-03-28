package weiboclock

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	// 导入statik生成的代码
	_ "cuitclock/statik"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/spf13/viper"
)

// 返回整点报时任务
func tollJob() cronweibo.WeiboJob {
	return cronweibo.WeiboJob{
		Name:     "toll",
		Schedule: "@hourly",
		Run:      tollRun,
	}
}

// 返回整点报时的文字和图片，用于创建job
func tollRun() (string, io.Reader) {
	// 生成文本内容
	now := Clock.Now()
	rand.Seed(now.Unix())
	mood := Moods[rand.Intn(len(Moods))]
	hour := now.Hour()
	oclock := hour
	// 12小时制处理
	if hour > 12 {
		oclock = hour - 12
	} else if hour == 0 {
		oclock = 12
	}
	words := strings.Repeat(Voices[rand.Intn(len(Voices))], oclock)
	text := fmt.Sprintf("%d点啦~ %s %s", oclock, mood, words)

	// 生成图片内容
	pic, err := PicReader(viper.GetString("weiboclock.pic_path"), hour)
	if err != nil {
		log.Println("[WARN] weiboclock toll error:", err)
		// 有error也不影响发送，获取图片失败就不发图片
	} else {
		if f, ok := pic.(*os.File); ok {
			defer f.Close()
		}
	}
	return text, pic
}
