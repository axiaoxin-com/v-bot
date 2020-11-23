// Package weiboclock 成信钟楼 on weibo
// 微博整点报时
package weiboclock

import (
	"github.com/axiaoxin-com/cronweibo"
	"github.com/axiaoxin-com/logging"
)

// WeiboClock 实例对象
type WeiboClock struct {
	cronWeibo *cronweibo.CronWeibo
}

// New 创建weiboclock实例
func New(cfg *cronweibo.Config) (*WeiboClock, error) {
	cw, err := cronweibo.New(cfg)
	if err != nil {
		return nil, err
	}

	// 创建实例
	clock := &WeiboClock{
		cronWeibo: cw,
	}

	return clock, nil
}

// Run 运行服务
func (clock *WeiboClock) Run() {
	// 初始化微博表情，失败不影响服务
	if count, err := clock.InitEmotions(); err != nil {
		logging.Warn(nil, "weiboclock InitEmotions error:"+err.Error())
	} else {
		logging.Debugf(nil, "weiboclock InitEmotions count:%d", count)
	}

	// 注册微博报时任务
	clock.cronWeibo.RegisterWeiboJobs(clock.tollJob())

	// 运行
	clock.cronWeibo.Start()
}
