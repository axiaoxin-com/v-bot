// Package weiboclock 成信钟楼 on weibo
// 微博整点报时
package weiboclock

import (
	"log"

	"github.com/axiaoxin-com/cronweibo"
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
		log.Println("[ERROR] weiboclock InitEmotions error", err)
	} else {
		log.Println("[DEBUG] weiboclock InitEmotions count:", count)
	}

	// 注册微博报时任务
	clock.cronWeibo.RegisterWeiboJobs(clock.tollJob())

	// 注册获取wttrin天气信息的普通定时任务
	// 整点请求wttrin响应大概率很慢或者会异常，导致报时延迟很大，提前五分钟获取天气保存在变量中，报时时直接从变量取值
	clock.cronWeibo.RegisterCronJobs(wttrinJob())

	// 运行
	clock.cronWeibo.Start()
}
