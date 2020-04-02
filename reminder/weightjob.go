package reminder

import (
	"fmt"
	"io"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/spf13/viper"
)

// 定时更新天气全局变量
func (r *Reminder) weightJob() cronweibo.WeiboJob {
	return cronweibo.WeiboJob{
		Name:     "weight",
		Schedule: viper.GetString("reminder.weight_schedule"),
		Run:      r.weightRun,
	}
}

// 生成天气信息
func (r *Reminder) weightRun() (string, io.Reader) {
	// 默认每月最后一天早上8点提醒称体重
	viper.SetDefault("reminder.wttrin_refresh_schedule", "0 0 8 L * ?")
	remindStr := r.RemindStr()
	text := fmt.Sprintf("该称体重了 %s", remindStr)
	return text, nil
}
