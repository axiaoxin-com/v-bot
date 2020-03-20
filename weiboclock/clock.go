// Package weiboclock 微博上的成信钟楼
// 通过cronserver实现微博报时
package weiboclock

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/axiaoxin/weibo"
	"github.com/pkg/errors"
)

// Clock 钟楼结构
type Clock struct {
	weibo          *weibo.Weibo
	token          *weibo.TokenResp
	tokenCreatedAt int64
	securityDomain string
}

// NewClock return clock object
func NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain string) (*Clock, error) {
	weibo := weibo.New(appkey, appsecret, username, passwd, redirecturi)
	if err := weibo.PCLogin(); err != nil {
		return nil, errors.Wrap(err, "weiboclock NewClock PCLogin error")
	}
	code, err := weibo.Authorize()
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock NewClock Authorize error")
	}
	token, err := weibo.AccessToken(code)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock NewClock AccessToken error")
	}
	return &Clock{
		weibo:          weibo,
		token:          token,
		tokenCreatedAt: time.Now().Unix(),
		securityDomain: securityDomain,
	}, nil
}

// OclockText 返回报时当前整点数和微博文本内容
func (c *Clock) OclockText() (int, string) {
	rand.Seed(time.Now().Unix())
	mood := Moods[rand.Intn(len(Moods))]
	oclock := time.Now().Hour()
	if oclock > 12 {
		oclock = oclock - 12
	}
	words := strings.Repeat(Words[rand.Intn(len(Words))], oclock)
	return oclock, fmt.Sprintf("%d点啦~ %s %s http://%s", oclock, mood, words, c.securityDomain)
}

// Toll 发送整点报时微博
// picPlan 使用的图片方案，图片路径
func (c *Clock) Toll(picPlan, picPath string) error {
	if err := c.UpdateToken(); err != nil {
		return errors.Wrap(err, "weiboclock Toll UpdateToken error")
	}
	oclock, text := c.OclockText()
	pic, err := PicReader(picPlan, picPath, oclock)
	if err != nil {
		log.Println("[WARN] weiboclock Toll error:", err)
		// 有error也不影响发送
	} else {
		if f, ok := pic.(*os.File); ok {
			defer f.Close()
		}
	}
	if err := c.weibo.StatusesShare(c.token.AccessToken, text, pic); err != nil {
		return errors.Wrap(err, "weiboclock Toll StatusesShare error")
	}
	return nil
}

// UpdateToken 检查access_token是否过去，过期则更新
func (c *Clock) UpdateToken() error {
	// 判断到当前时间为止token已存在时间是否已大于其过期时间
	age := time.Now().Unix() - c.tokenCreatedAt
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
