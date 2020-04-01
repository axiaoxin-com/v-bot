package weiboclock

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/schollz/progressbar/v2"
)

// ProgressBar 返回静态进度条
func ProgressBar(theme progressbar.Theme, width, total, current int) string {
	buf := strings.Builder{}
	bar := progressbar.NewOptions(
		total,
		progressbar.OptionSetTheme(theme),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetWidth(width),
		progressbar.OptionSetWriter(&buf),
	)

	if err := bar.Set(current); err != nil {
		log.Println("[ERROR] progressbar Set error", err)
	}
	return strings.TrimSpace(buf.String())
}

// DayProgressBar 今日使用进度
func DayProgressBar(t time.Time) string {
	// 使用day的时间戳作为seed，一天内使用相同主题
	ts := t.Unix()
	ts = ts - ts%(60*60*24)
	rand.Seed(ts)
	saucerAndPaddings := [][]string{
		{"░", "▒"},
		{"⬛️", "⬜️"},
		{"❌", "⭕️"},
		{"⚫️", "⚪️"},
		{"🖤", "🤍"},
		{"🤍", "❤️"},
	}
	saucerAndPadding := saucerAndPaddings[rand.Intn(len(saucerAndPaddings))]
	theme := progressbar.Theme{Saucer: saucerAndPadding[0], SaucerHead: "", SaucerPadding: saucerAndPadding[1], BarStart: "", BarEnd: ""}

	hour := t.Hour()
	if hour == 0 {
		hour = 24
	}
	bar := ProgressBar(theme, 10, 24, hour)
	// 替换 [hour:24] 为 [hour/24]
	bar = strings.Replace(bar, ":", "/", 1)
	return bar
}
