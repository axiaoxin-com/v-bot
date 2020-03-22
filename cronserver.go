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
var weiboDebugClock *weiboclock.Clock
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

	if viper.GetBool("server.debug") {
		testUsername := viper.GetString("weibo.test_username")
		testPasswd := viper.GetString("weibo.test_passwd")
		testAuthCode := viper.GetString("weibo.test_auth_code")
		weiboDebugClock, err = weiboclock.NewClock(appkey, appsecret, testUsername, testPasswd, redirecturi, securityDomain, testAuthCode)
		if err != nil {
			log.Println("[ERROR] cronserver init weibo debug clock error:", err)
		} else {
			log.Println("[INFO] cronserver inited weiboDebugClock.")
		}
	}
	return nil
}

// 执行微博报时
func doToll(useDebugClock bool) string {
	if weiboClock == nil {
		log.Println("[WARN] cronserver tollJob find weiboClock is nil, try to initWeiboClock...")
		if err := initWeiboClock(); err != nil {
			return ""
		}
	}

	vbClock := weiboClock
	if useDebugClock {
		vbClock = weiboDebugClock
	}
	picPath := viper.GetString("weiboclock.pic_path")
	resp, err := vbClock.Toll(picPath)
	if err != nil {
		log.Println("[ERROR] cronserver tollJob Toll error:", err)
	}
	homeURL := "http://weibo.com/" + resp.User.ProfileURL
	log.Println("[INFO] cronserver doToll complete.", homeURL)
	return homeURL
}

// 微博报时任务
func tollJob() {
	doToll(false)
}

func runCronServer() {
	initWeiboClock()
	cronLocation := viper.GetString("cron.location")
	location, err := time.LoadLocation(cronLocation)
	if err != nil {
		log.Fatal("[FATAL] cronserver load location error:", err)
	}
	c := cron.NewWithLocation(location)
	log.Println("[INFO] cronserver adding jobs...")
	if ringJobSchedule := viper.GetString("cron.toll_job"); ringJobSchedule != "" {
		if err := c.AddFunc(ringJobSchedule, tollJob); err != nil {
			log.Println("[ERROR] cronserver add tollJob error:", err)
		} else {
			log.Println("[INFO] cronserver added tollJob as", ringJobSchedule)
		}
	}
	log.Println("[INFO] cronserver is running with location", location)
	c.Start()
	defer c.Stop()
	select {}
}
