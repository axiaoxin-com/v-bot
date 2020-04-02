package reminder

import (
	"fmt"
	"io"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/spf13/viper"
)

func init() {
	// 默认每月第1天早上8点提醒称体重
	viper.SetDefault("reminder.weight_schedule", "0 0 8 1 * *")
}

// 定时提醒称体重
func (r *Reminder) weightJob() cronweibo.WeiboJob {
	return cronweibo.WeiboJob{
		Name:     "weight",
		Schedule: viper.GetString("reminder.weight_schedule"),
		Run:      r.weightRun,
	}
}

func (r *Reminder) weightRun() (string, io.Reader) {
	remindStr := r.RemindStr()
	text := fmt.Sprintf("该称体重了 %s", remindStr)
	return text, nil
}
