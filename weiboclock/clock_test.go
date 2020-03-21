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

	clock, err := NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain)
	if err != nil {
		t.Fatal(err)
	}
	// test local plan
	picPlan := "local"
	picPath := "../pictures/weibo"
	if err := clock.Toll(picPlan, picPath); err != nil {
		t.Error(err)
	}
}
