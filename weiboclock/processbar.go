package weiboclock

import (
	"log"
	"strings"

	"github.com/schollz/progressbar/v2"
)

// ProgressBar 返回静态进度条
func ProgressBar(width, total, current int) string {
	buf := strings.Builder{}
	theme := progressbar.Theme{Saucer: "🤍", SaucerHead: "", SaucerPadding: "❤️", BarStart: "", BarEnd: ""}
	// theme := progressbar.Theme{Saucer: "░", SaucerHead: "", SaucerPadding: "▒", BarStart: "", BarEnd: ""}
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
func DayProgressBar(hour int) string {
	if hour == 0 {
		hour = 24
	}
	bar := ProgressBar(10, 24, hour)
	// 替换 [hour:24] 为 [hour/24]
	bar = strings.Replace(bar, ":", "/", 1)
	return bar
}
