package weibo

import (
	"cuitclock/config"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestStatusesShare(t *testing.T) {
	config.InitConfig("..")
	tAppkey := viper.GetString("weibo.app_key")
	tAppsecret := viper.GetString("weibo.app_secret")
	tUsername := viper.GetString("weibo.test_username")
	tPasswd := viper.GetString("weibo.test_passwd")
	tRedirecturi := viper.GetString("weibo.redirect_uri")
	tSecurityDomain := viper.GetString("weibo.security_domain")
	weibo := NewWeibo(tAppkey, tAppsecret, tUsername, tPasswd, tRedirecturi)
	if err := weibo.PCLogin(); err != nil {
		t.Fatal(err)
	}
	code, err := weibo.AuthCode()
	if err != nil {
		t.Fatal(err)
	}
	token, err := weibo.AccessToken(code)
	if err != nil {
		t.Fatal(err)
	}
	status := fmt.Sprintf("%s http://%s", time.Now().Format("2006-01-02 15:04:05"), tSecurityDomain)
	if err := weibo.StatusesShare(token.AccessToken, status, nil); err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	pic, err := os.Open("./pic.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer pic.Close()
	if err := weibo.StatusesShare(token.AccessToken, status, pic); err != nil {
		t.Error(err)
	}
}
