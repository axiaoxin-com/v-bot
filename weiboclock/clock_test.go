package weiboclock

import (
	"cuitclock/config"
	"testing"

	"github.com/spf13/viper"
)

func TestToll(t *testing.T) {
	config.InitConfig("..")
	tAppkey := viper.GetString("weibo.app_key")
	tAppsecret := viper.GetString("weibo.app_secret")
	tUsername := viper.GetString("weibo.test_username")
	tPasswd := viper.GetString("weibo.test_passwd")
	tRedirecturi := viper.GetString("weibo.redirect_uri")
	tSecurityDomain := viper.GetString("weibo.security_domain")

	clock, err := NewClock(tAppkey, tAppsecret, tUsername, tPasswd, tRedirecturi, tSecurityDomain)
	if err != nil {
		t.Fatal(err)
	}
	if err := clock.Toll(); err != nil {
		t.Error(err)
	}
}
