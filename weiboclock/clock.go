package weiboclock

import (
	"cuitclock/weibo"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	weibo          *weibo.Weibo
	token          *weibo.TokenResp
	tokenCreatedAt int64
	securityDomain string
}

// NewClock return clock object
func NewClock(appkey, appsecret, username, passwd, redirecturi, securityDomain string) (*Clock, error) {
	weibo := weibo.NewWeibo(appkey, appsecret, username, passwd, redirecturi)
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
		token:          token,
		tokenCreatedAt: time.Now().Unix(),
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

// Toll 发送整点报时微博
func (c *Clock) Toll() error {
	if err := c.UpdateToken(); err != nil {
		return errors.Wrap(err, "clock Toll UpdateToken error")
	}
	text := c.OclockText()
	if err := c.weibo.StatusesShare(c.token.AccessToken, text, nil); err != nil {
		return errors.Wrap(err, "clock Toll StatusesShare error")
	}
	return nil
}

// UpdateToken 检查access_token是否过去，过期则更新
func (c *Clock) UpdateToken() error {
	// 判断到当前时间为止token已存在时间是否已大于其过期时间
	age := time.Now().Unix() - c.tokenCreatedAt
	// 过期则更新token
	if age >= c.token.ExpiresIn {
		log.Println("[INFO] clock token will expire, let set a new token")
		if err := c.weibo.PCLogin(); err != nil {
			return errors.Wrap(err, "clock UpdateToken PCLogin error")
		}
		code, err := c.weibo.AuthCode()
		if err != nil {
			return errors.Wrap(err, "clock UpdateToken AuthCode error")
		}
		token, err := c.weibo.AccessToken(code)
		if err != nil {
			return errors.Wrap(err, "clock UpdateToken AccessToken error")
		}
		c.token = token
	}
	log.Println("[INFO] clock check token age =", age)
	return nil
}
