package main

import (
	"log"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

func ringJob() {
	clock, err := NewClock()
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
	cronLocation := viper.GetString("cron.location")
	if cronLocation == "" {
		cronLocation = "Asia/Shanghai"
	}
	location, err := time.LoadLocation(cronLocation)
	if err != nil {
		log.Fatal("load location error:", err)
	}
	c := cron.NewWithLocation(location)
	if ringJobSchedule := viper.GetString("cron.ring_job"); ringJobSchedule != "" {
		c.AddFunc(ringJobSchedule, ringJob)
		log.Println("add ring job as", ringJobSchedule)
	}
	c.Start()
	defer c.Stop()
	select {}
}
