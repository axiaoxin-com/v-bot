package weiboclock

import (
	"cuitclock/config"
	"testing"

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
	picPlan := viper.GetString("weiboclock.pic_plan")
	picPath := viper.GetString("weiboclock.pic_path")

	clock, err := NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain)
	if err != nil {
		t.Fatal(err)
	}
	if err := clock.Toll(picPlan, picPath); err != nil {
		t.Error(err)
	}
}
