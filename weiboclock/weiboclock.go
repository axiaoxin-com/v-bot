// Package weiboclock 成信钟楼 on weibo
// 微博整点报时
package weiboclock

import (
	"log"

	"github.com/axiaoxin-com/cronweibo"
)

// Clock weiboclock对象
var Clock *cronweibo.CronWeibo

// Run 运行weiboclock
func Run(cfg *cronweibo.Config) {
	var err error
	Clock, err = cronweibo.New(cfg)
	if err != nil {
		log.Fatalln("[FATAL] New cronweibo error:", err)
	}

	// 注册报时任务
	Clock.RegisterWeiboJobs(tollJob())

	// 运行
	Clock.Start()
}
