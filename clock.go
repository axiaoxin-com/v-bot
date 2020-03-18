package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// Words 微博文字
	Words = []string{
		"biu! ", "ha! ", "mua~ ", "dang! ", "ho! ", "hei! ", "pia! ",
		"he! ", "mia~ ", "ao~ ", "do~ ", "ga~ ", "bi! ", "ba~ ", "ma~ ",
		"pa! ", "no~ ", "bomb! ", "yeah~ ", "ka! ", "la~ ", "da! ", "OA! ",
	}

	// Moods 微博表情
	Moods = []string{
		"(°ー°〃)", "_(:з」∠)_ ", "o(*≧▽≦)ツ┏━┓", "๑乛◡乛๑ ", "(σ‘・д･)σ",
		"( ＿ ＿)ノ｜", "┑(￣Д ￣)┍", "[放电抛媚][霹雳][霹雳][霹雳][被电]", "(＃°Д°)",
		"(-ω- )", "(′・ω・`)", "( ^ω^)", "乀(ˉεˉ乀)", "\\(╯-╰)/", "\\(▔▽▔)/",
		"^(oo)^", "(O^~^O)", "（╯＾╰）", "[呵呵]", "[哈哈]", "[生病]", "[委屈]", "[泪]",
		"[衰]", "[嘘]", "[悲伤]", "[怒骂]", "[伤心]", "[打哈欠]", "[走你]", "[江南style]",
		"[不想上班]", "[笑哈哈]", "[得意地笑]", "[泪流满面]", "[纠结]", "[抠鼻屎]", "[求关注]",
		"[奥特曼]", "[瞧瞧]", "[嘻嘻]", "[可爱]", "[可怜]", "[挖鼻屎]", "[吃惊]", "[害羞]",
		"[挤眼]", "[闭嘴]", "[鄙视]", "[爱你]", "[偷笑]", "[亲亲]", "[太开心]", "[懒得理你]",
		"[右哼哼]", "[左哼哼]", "[吐]", "[抱抱]", "[怒]", "[疑问]", "[馋嘴]", "[拜拜]",
		"[思考]", "[汗]", "[困]", "[睡觉]", "[甩甩手]", "[失望]", "[酷]", "[花心]", "[哼]",
		"[鼓掌]", "[拍手]", "[抓狂]", "[黑线]", "[阴险]", "[心]", "[偷乐]", "[转发]",
		"[好爱哦]", "[蜡烛]", "[羞嗒嗒]", "[大南瓜]", "[立志青年]", "[困死了]", "[带感]",
		"[崩溃]", "[好囧]", "[别烦我]", "[din害羞]", "[din吃]", "[lxhx喵喵]", "[g思考]",
		"[lm天然呆]", "[bed凌乱]", "[c捂脸]", "[乐乐]", "[ali踩]", "[冒个泡]", "[吵闹]",
		"[眨眨眼]",
	}
)

// Clock 钟楼结构
type Clock struct {
	weibo          *Weibo
	accessToken    string
	securityDomain string
}

// NewClock return clock object
func NewClock() (*Clock, error) {
	appkey := viper.GetString("weibo.app_key")
	appsecret := viper.GetString("weibo.app_secret")
	username := viper.GetString("weibo.username")
	passwd := viper.GetString("weibo.passwd")
	redirecturi := viper.GetString("weibo.redirect_uri")
	securityDomain := viper.GetString("weibo.security_domain")
	weibo := NewWeibo(appkey, appsecret, username, passwd, redirecturi)
	if err := weibo.PCLogin(); err != nil {
		return nil, errors.Wrap(err, "clock NewClock PCLogin error")
	}
	code, err := weibo.AuthCode()
	if err != nil {
		return nil, errors.Wrap(err, "clock NewClock AuthCode error")
	}
	token, err := weibo.AccessToken(code)
	if err != nil {
		return nil, errors.Wrap(err, "clock NewClock AccessToken error")
	}
	return &Clock{
		weibo:          weibo,
		accessToken:    token.AccessToken,
		securityDomain: securityDomain,
	}, nil
}

// OclockText 返回报时文本内容
func (c *Clock) OclockText() string {
	rand.Seed(time.Now().Unix())
	mood := Moods[rand.Intn(len(Moods))]
	oclock := time.Now().Hour()
	if oclock > 12 {
		oclock = oclock - 12
	}
	words := strings.Repeat(Words[rand.Intn(len(Words))], oclock)
	return fmt.Sprintf("%d点啦~ %s %s \nhttp://%s", oclock, mood, words, c.securityDomain)
}

// Ring 整点报时
func (c *Clock) Ring() error {
	text := c.OclockText()
	return c.weibo.StatusesShare(c.accessToken, text, nil)
}
