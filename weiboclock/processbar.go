package weiboclock

import (
	"log"
	"strings"

	"github.com/schollz/progressbar/v2"
)

// ProgressBar 返回静态进度条
func ProgressBar(width, total, current int) string {
	buf := strings.Builder{}
	bar := progressbar.NewOptions(
		total,
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "▓", SaucerHead: "", SaucerPadding: "░", BarStart: "", BarEnd: ""}),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetWidth(width),
		progressbar.OptionSetWriter(&buf),
		progressbar.OptionSetDescription(""),
		progressbar.OptionSetRenderBlankState(true),
	)

	if err := bar.Set(current); err != nil {
		log.Println("[ERROR] progressbar Set error", err)
	}
	return strings.TrimSpace(buf.String())
}
