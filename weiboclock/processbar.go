package weiboclock

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/schollz/progressbar/v2"
)

// ProgressBar 返回静态进度条
func ProgressBar(width, total, current int) string {
	buf := strings.Builder{}
	saucerAndPaddings := [][]string{
		{"░", "▒"},
		{"⬛️", "⬜️"},
		{"❌", "⭕️"},
		{"⚫️", "⚪️"},
		{"🖤", "🤍"},
		{"🤍", "❤️"},
	}
	rand.Seed(time.Now().Unix())
	saucerAndPadding := saucerAndPaddings[rand.Intn(len(saucerAndPaddings))]

	theme := progressbar.Theme{Saucer: saucerAndPadding[0], SaucerHead: "", SaucerPadding: saucerAndPadding[1], BarStart: "", BarEnd: ""}
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
