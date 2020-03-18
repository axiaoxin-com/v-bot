package main

import (
	"log"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var clock *Clock
var err error

func ringJob() {
	if err != nil {
		log.Println("NewClock error:", err)
		return
	}
	if err := clock.Ring(); err != nil {
		log.Println("ringJob Ring error:", err)
	}
	log.Println("ringJob complete.")
}

func runCronServer() {
	appkey := viper.GetString("weibo.app_key")
	appsecret := viper.GetString("weibo.app_secret")
	username := viper.GetString("weibo.username")
	passwd := viper.GetString("weibo.passwd")
	redirecturi := viper.GetString("weibo.redirect_uri")
	securityDomain := viper.GetString("weibo.security_domain")
	clock, err = NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain)
	cronLocation := viper.GetString("cron.location")
	if cronLocation == "" {
		cronLocation = "Asia/Shanghai"
	}
	location, err := time.LoadLocation(cronLocation)
	if err != nil {
		log.Fatal("load location error:", err)
	}
	log.Println("run cron server with location", location)
	c := cron.NewWithLocation(location)
	log.Println("adding jobs for cron server...")
	if ringJobSchedule := viper.GetString("cron.ring_job"); ringJobSchedule != "" {
		if err := c.AddFunc(ringJobSchedule, ringJob); err != nil {
			log.Println("add ring job error:", err)
		} else {
			log.Println("added ring job as", ringJobSchedule)
		}
	}
	c.Start()
	defer c.Stop()
	select {}
}
