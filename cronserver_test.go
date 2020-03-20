package main

import (
	"cuitclock/config"
	"testing"

	"github.com/spf13/viper"
)

func TestTollJob(t *testing.T) {
	config.InitConfig()
	viper.Set("weibo.username", viper.GetString("weibo.test_username"))
	viper.Set("weibo.passwd", viper.GetString("weibo.test_passwd"))
	tollJob()
}
