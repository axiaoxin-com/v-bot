// Package weiboclock 微博上的成信钟楼
// 通过cronserver实现微博报时
package weiboclock

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	// 导入statik生成的代码
	_ "cuitclock/statik"

	"github.com/axiaoxin-com/weibo"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/viper"
)

// StatikFS statik obj
var StatikFS http.FileSystem

func init() {
	var err error
	StatikFS, err = fs.New()
	if err != nil {
		log.Fatalln("cuitclock init StatikFS error:", err)
	}
}

// Clock 钟楼结构
type Clock struct {
	weibo          *weibo.Weibo
	token          *weibo.TokenResp
	tokenCreatedAt int64
	securityDomain string
}

// NewClock return clock object
func NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain, authCode string) (*Clock, error) {
	weibo := weibo.New(appkey, appsecret, username, passwd, redirecturi)
	if authCode == "" {
		if err := weibo.PCLogin(); err != nil {
			return nil, errors.Wrap(err, "weiboclock NewClock PCLogin error")
		}
		code, err := weibo.Authorize()
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock NewClock Authorize error")
		}
		authCode = code
	}
	token, err := weibo.AccessToken(authCode)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock NewClock AccessToken error")
	}
	log.Println("[DEBUG] weiboclock NewClock code:", authCode, " token:", token.AccessToken)
	return &Clock{
		weibo:          weibo,
		token:          token,
		tokenCreatedAt: LocationNow().Unix(),
		securityDomain: securityDomain,
	}, nil
}

// OclockText 返回报时当前24小时制hour数和微博文本内容
func (c *Clock) OclockText() (int, string) {
	now := LocationNow()
	rand.Seed(now.Unix())
	mood := Moods[rand.Intn(len(Moods))]
	hour := now.Hour()
	oclock := hour
	// 12小时制处理
	if hour > 12 {
		oclock = hour - 12
	} else if hour == 0 {
		oclock = 12
	}
	words := strings.Repeat(Voices[rand.Intn(len(Voices))], oclock)
	return hour, fmt.Sprintf("%d点啦~ %s %s http://%s", oclock, mood, words, c.securityDomain)
}

// Toll 发送整点报时微博
// picPath 指定图片路径
func (c *Clock) Toll(picPath string) (*weibo.StatusesShareResp, error) {
	if err := c.UpdateToken(); err != nil {
		return nil, errors.Wrap(err, "weiboclock Toll UpdateToken error")
	}
	hour, text := c.OclockText()
	pic, err := PicReader(picPath, hour)
	if err != nil {
		log.Println("[WARN] weiboclock Toll error:", err)
		// 有error也不影响发送
	} else {
		if f, ok := pic.(*os.File); ok {
			defer f.Close()
		}
	}
	resp, err := c.weibo.StatusesShare(c.token.AccessToken, text, pic)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock Toll StatusesShare error")
	}
	return resp, nil
}

// UpdateToken 检查access_token是否过去，过期则更新
func (c *Clock) UpdateToken() error {
	// 判断到当前时间为止token已存在时间是否已大于其过期时间
	age := LocationNow().Unix() - c.tokenCreatedAt
	// 过期则更新token
	if age >= c.token.ExpiresIn {
		log.Println("[INFO] weiboclock token will expire, let set a new token")
		if err := c.weibo.PCLogin(); err != nil {
			return errors.Wrap(err, "weiboclock UpdateToken PCLogin error")
		}
		code, err := c.weibo.Authorize()
		if err != nil {
			return errors.Wrap(err, "weiboclock UpdateToken Authorize error")
		}
		token, err := c.weibo.AccessToken(code)
		if err != nil {
			return errors.Wrap(err, "weiboclock UpdateToken AccessToken error")
		}
		c.token = token
	}
	log.Printf("[INFO] weiboclock check token age=%d, ExpiresIn=%d", age, c.token.ExpiresIn)
	return nil
}

// LocationNow 返回配置中cron.location时区的当前时间
func LocationNow() time.Time {
	loc, err := time.LoadLocation(viper.GetString("cron.location"))
	if err != nil {
		log.Println("[ERROR] weiboclock LocationNow LoadLocation error:", err)
		return time.Now()
	}
	now := time.Now().In(loc)
	return now
}
