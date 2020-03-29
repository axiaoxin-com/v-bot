package main

import (
	"cuitclock/config"
	"cuitclock/weiboclock"
	"log"
	"time"

	"github.com/axiaoxin-com/chaojiying"
	"github.com/axiaoxin-com/cronweibo"
	"github.com/axiaoxin-com/weibo"
	"github.com/spf13/viper"
)

// 运行微博上的成信钟楼
func runWeiboClock() {
	// 初始化 weiboclock 的配置
	location, err := time.LoadLocation(viper.GetString("weiboclock.location"))
	if err != nil {
		log.Fatalln("[FATAL] Load location error:", err)
	}
	username := viper.GetString("weibo.username")
	passwd := viper.GetString("weibo.passwd")
	tusername := viper.GetString("weibo.test_username")
	tpasswd := viper.GetString("weibo.test_passwd")
	if tusername != "" && tpasswd != "" {
		username = tusername
		passwd = tpasswd
	}
	wcCfg := &cronweibo.Config{
		WeiboAppkey:        viper.GetString("weibo.app_key"),
		WeiboAppsecret:     viper.GetString("weibo.app_secret"),
		WeiboUsername:      username,
		WeiboPasswd:        passwd,
		WeiboRedirecturi:   viper.GetString("weibo.redirect_uri"),
		WeiboSecurityURL:   viper.GetString("weibo.security_url"),
		WeiboPinCrackFuncs: []weibo.CrackPinFunc{},
		HTTPServerAddr:     viper.GetString("weiboclock.webapi_addr"),
		BasicAuthUsername:  viper.GetString("weiboclock.basic_auth_username"),
		BasicAuthPasswd:    viper.GetString("weiboclock.basic_auth_passwd"),
		Location:           location,
	}
	// 使用超级鹰破解验证码
	// 初始化超级鹰客户端
	accountsJSONPath := viper.GetString("chaojiying.accounts_json_path")
	if accountsJSONPath != "" {
		accounts, err := chaojiying.LoadAccountsFromJSONFile(accountsJSONPath)
		if err != nil {
			log.Fatal("[FATAL] Load chaojiying accounts error:", err)
		}
		cracker, err := chaojiying.New(accounts)
		if err != nil {
			log.Fatal("[FATAL] New chaojiying cracker error:", err)
		}
		wcCfg.WeiboPinCrackFuncs = []weibo.CrackPinFunc{cracker.Cr4ck}
	}

	// 运行weiboclock
	weiboClock, err := weiboclock.New(wcCfg)
	if err != nil {
		log.Fatal(err)
	}
	weiboClock.Run()
}

func main() {
	config.InitConfig()
	log.Println("[INFO] cuitclock inited config.")
	runWeiboClock()
}
