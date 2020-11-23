package main

import (
	"fmt"
	"time"
	"v-bot/config"
	"v-bot/reminder"
	"v-bot/weiboclock"

	"github.com/axiaoxin-com/chaojiying"
	"github.com/axiaoxin-com/cronweibo"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/weibo"
	"github.com/spf13/viper"
)

const (
	// DefaultTimezone 时区
	DefaultTimezone = "Asia/Shanghai"
)

// 运行微博上的成信钟楼
func runWeiboClock(cracker *chaojiying.Client) {
	// 初始化 weiboclock 的配置
	timezone := viper.GetString("weiboclock.timezone")
	if timezone == "" {
		timezone = DefaultTimezone
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	username := viper.GetString("weiboclock.username")
	passwd := viper.GetString("weiboclock.passwd")
	tusername := viper.GetString("weiboclock.test_username")
	tpasswd := viper.GetString("weiboclock.test_passwd")
	if tusername != "" && tpasswd != "" {
		username = tusername
		passwd = tpasswd
	}
	wcCfg := &cronweibo.Config{
		AppName:            "WeiboClock",
		WeiboAppkey:        viper.GetString("weiboclock.app_key"),
		WeiboAppsecret:     viper.GetString("weiboclock.app_secret"),
		WeiboUsername:      username,
		WeiboPasswd:        passwd,
		WeiboRedirecturi:   viper.GetString("weiboclock.redirect_uri"),
		WeiboSecurityURL:   viper.GetString("weiboclock.security_url"),
		WeiboPinCrackFuncs: []weibo.CrackPinFunc{},
		HTTPServerAddr:     viper.GetString("weiboclock.webapi_addr"),
		BasicAuthUsername:  viper.GetString("weiboclock.basic_auth_username"),
		BasicAuthPasswd:    viper.GetString("weiboclock.basic_auth_passwd"),
		Location:           location,
	}
	if cracker != nil {
		wcCfg.WeiboPinCrackFuncs = []weibo.CrackPinFunc{cracker.Cr4ck}
	}

	// 运行weiboclock
	weiboClock, err := weiboclock.New(wcCfg)
	if err != nil {
		panic(err)
	}
	weiboClock.Run()
}

// 运行微博提醒事项
func runReminder(cracker *chaojiying.Client) {
	// 初始化 weiboclock 的配置
	timezone := viper.GetString("reminder.timezone")
	if timezone == "" {
		timezone = DefaultTimezone
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	username := viper.GetString("reminder.username")
	passwd := viper.GetString("reminder.passwd")
	tusername := viper.GetString("reminder.test_username")
	tpasswd := viper.GetString("reminder.test_passwd")
	if tusername != "" && tpasswd != "" {
		username = tusername
		passwd = tpasswd
	}
	wcCfg := &cronweibo.Config{
		AppName:            "Reminder",
		WeiboAppkey:        viper.GetString("reminder.app_key"),
		WeiboAppsecret:     viper.GetString("reminder.app_secret"),
		WeiboUsername:      username,
		WeiboPasswd:        passwd,
		WeiboRedirecturi:   viper.GetString("reminder.redirect_uri"),
		WeiboSecurityURL:   viper.GetString("reminder.security_url"),
		WeiboPinCrackFuncs: []weibo.CrackPinFunc{},
		HTTPServerAddr:     viper.GetString("reminder.webapi_addr"),
		BasicAuthUsername:  viper.GetString("reminder.basic_auth_username"),
		BasicAuthPasswd:    viper.GetString("reminder.basic_auth_passwd"),
		Location:           location,
	}
	if cracker != nil {
		wcCfg.WeiboPinCrackFuncs = []weibo.CrackPinFunc{cracker.Cr4ck}
	}
	// 运行reminder
	r, err := reminder.New(wcCfg)
	if err != nil {
		panic(err)
	}
	r.Run()
}

// 获取超级鹰客户端
func cracker() *chaojiying.Client {
	// 使用超级鹰破解验证码
	// 初始化超级鹰客户端
	accountsJSONPath := viper.GetString("chaojiying.accounts_json_path")
	if accountsJSONPath != "" {
		accounts, err := chaojiying.LoadAccountsFromJSONFile(accountsJSONPath)
		if err != nil {
			logging.Error(nil, "Load chaojiying accounts error:"+err.Error())
		}
		cracker, err := chaojiying.New(accounts)
		if err != nil {
			logging.Error(nil, "New chaojiying cracker error:"+err.Error())
		}
		return cracker
	}
	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			select {}
		}
	}()
	if err := config.InitConfig(); err != nil {
		panic("InitConfig:" + err.Error())
	}
	c := cracker()
	go runWeiboClock(c)
	go runReminder(c)
	select {}
}
