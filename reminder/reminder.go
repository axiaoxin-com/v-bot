// Package reminder 微博提醒事项
package reminder

import "github.com/axiaoxin-com/cronweibo"

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

	// 运行
	r.cronWeibo.Start()
}
