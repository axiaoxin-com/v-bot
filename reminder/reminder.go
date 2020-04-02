// Package reminder 微博提醒事项
package reminder

import (
	"strings"

	"github.com/axiaoxin-com/cronweibo"
	"github.com/spf13/viper"
)

// Reminder 实例对象
type Reminder struct {
	cronWeibo *cronweibo.CronWeibo
}

// New 创建weiboclock实例
func New(cfg *cronweibo.Config) (*Reminder, error) {
	cw, err := cronweibo.New(cfg)
	if err != nil {
		return nil, err
	}

	// 创建实例
	r := &Reminder{
		cronWeibo: cw,
	}

	return r, nil
}

// Run 运行服务
func (r *Reminder) Run() {
	// 注册天气变化提醒任务
	r.cronWeibo.RegisterWeiboJobs(r.weatherJob())
	// 注册称体重提醒任务
	r.cronWeibo.RegisterWeiboJobs(r.weightJob())

	// 运行
	r.cronWeibo.Start()
}

// RemindStr 获取微博提醒昵称列表，空格分隔
func (r *Reminder) RemindStr() string {
	nicknameList := strings.Fields(viper.GetString("reminder.remind_list"))
	remindList := []string{}
	for _, nickname := range nicknameList {
		if !strings.HasPrefix(nickname, "@") {
			remindList = append(remindList, "@"+nickname)
		} else {
			remindList = append(remindList, nickname)
		}
	}
	return strings.Join(remindList, " ")
}
