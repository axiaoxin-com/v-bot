// 定时任务

package main

import (
	"cuitclock/weiboclock"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var weiboClock *weiboclock.Clock
var err error

func initWeiboClock() error {
	appkey := viper.GetString("weibo.app_key")
	appsecret := viper.GetString("weibo.app_secret")
	username := viper.GetString("weibo.username")
	passwd := viper.GetString("weibo.passwd")
	redirecturi := viper.GetString("weibo.redirect_uri")
	securityDomain := viper.GetString("weibo.security_domain")
	authCode := viper.GetString("weibo.auth_code")
	authURL := fmt.Sprintf("https://api.weibo.com/oauth2/authorize?redirect_uri=%s&response_type=code&client_id=%s", redirecturi, appkey)
	log.Println("[INFO] authorize url:", authURL)
	weiboClock, err = weiboclock.NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain, authCode)
	if err != nil {
		log.Println("[ERROR] cronserver init weibo clock error:", err)
		return errors.Wrap(err, "cronserver initClock error")
	}
	log.Println("[INFO] cronserver inited weiboClock.")
	return nil
}

// 微博报时任务
func tollJob() {
	if weiboClock == nil {
		log.Println("[WARN] cronserver tollJob find weiboClock is nil, try to initWeiboClock...")
		if err := initWeiboClock(); err != nil {
			return
		}
	}
	picPath := viper.GetString("weiboclock.pic_path")
	if err := weiboClock.Toll(picPath); err != nil {
		log.Println("[ERROR] cronserver tollJob Toll error:", err)
	}
	log.Println("[INFO] cronserver tollJob complete.")
}

func runCronServer() {
	initWeiboClock()
	cronLocation := viper.GetString("cron.location")
	if cronLocation == "" {
		cronLocation = "Asia/Shanghai"
	}
	location, err := time.LoadLocation(cronLocation)
	if err != nil {
		log.Fatal("[FATAL] cronserver load location error:", err)
	}
	log.Println("[INFO] cronserver running with location", location)
	c := cron.NewWithLocation(location)
	log.Println("[INFO] cronserver adding jobs...")
	if ringJobSchedule := viper.GetString("cron.toll_job"); ringJobSchedule != "" {
		if err := c.AddFunc(ringJobSchedule, tollJob); err != nil {
			log.Println("[ERROR] cronserver add tollJob error:", err)
		} else {
			log.Println("[INFO] cronserver added tollJob as", ringJobSchedule)
		}
	}
	c.Start()
	defer c.Stop()
	select {}
}
