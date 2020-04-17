package reminder

import (
	"fmt"
	"io"
	"v-bot/progressbar"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/spf13/viper"
)

func init() {
	// 默认每月第一天早上 9 点提醒今年时间进度
	viper.SetDefault("reminder.yearbar_schedule", "0 9 1 * *")
}

// 定时提醒今年时间进度条
func (r *Reminder) yearbarJob() cronweibo.WeiboJob {
	return cronweibo.WeiboJob{
		Name:     "yearbar",
		Schedule: viper.GetString("reminder.yearbar_schedule"),
		Run:      r.yearbarRun,
	}
}

// 生成 yearbar
func (r *Reminder) yearbarRun() (string, io.Reader) {
	remindStr := r.RemindStr()
	now := r.cronWeibo.Now()
	bar := progressbar.YearProgressBar(now)
	text := fmt.Sprintf("你的%d使用进度:\n\n%s\n\n——%s", now.Year(), bar, remindStr)
	return text, nil
}
