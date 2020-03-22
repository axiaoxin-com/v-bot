package weiboclock

import (
	"cuitclock/config"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestToll(t *testing.T) {
	config.InitConfig("..")
	appkey := viper.GetString("weibo.app_key")
	appsecret := viper.GetString("weibo.app_secret")
	username := viper.GetString("weibo.test_username")
	passwd := viper.GetString("weibo.test_passwd")
	redirecturi := viper.GetString("weibo.redirect_uri")
	securityDomain := viper.GetString("weibo.security_domain")
	authCode := viper.GetString("weibo.auth_code")

	clock, err := NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain, authCode)
	if err != nil {
		t.Fatal(err)
	}
	// 测试指定图片
	picPath := "../assets/weibo"
	if err := clock.Toll(picPath); err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	// 测试内置图片
	picPath = "default"
	if err := clock.Toll(picPath); err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	// 测试无图片
	picPath = ""
	if err := clock.Toll(picPath); err != nil {
		t.Error(err)
	}
}
