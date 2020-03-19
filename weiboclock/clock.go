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
		"(°ー°〃)", "_(:з」∠)_ ", "o(*≧▽≦)ツ┏━┓", "๑乛◡乛๑ ", "(σ‘・д･)σ", "( ＿ ＿)ノ｜", "┑(￣Д ￣)┍",
		"[放电抛媚][霹雳][霹雳][霹雳][被电]", "(＃°Д°)", "(-ω- )", "(′・ω・`)", "( ^ω^)", "乀(ˉεˉ乀)",
		"\\(╯-╰)/", "\\(▔▽▔)/", "^(oo)^", "(O^~^O)", "（╯＾╰）",
		"[武汉加油]", "[点亮平安灯]", "[炸鸡腿]", "[中国赞]", "[锦鲤]", "[抱抱]", "[摊手]", "[跪了]", "[酸]",
		"[花木兰]", "[BB8]", "[冰雪奇缘艾莎]", "[微笑]", "[可爱]", "[太开心]", "[鼓掌]", "[嘻嘻]", "[哈哈]",
		"[笑cry]", "[挤眼]", "[馋嘴]", "[黑线]", "[汗]", "[挖鼻]", "[哼]", "[怒]", "[委屈]", "[可怜]", "[失望]",
		"[悲伤]", "[泪]", "[允悲]", "[害羞]", "[污]", "[爱你]", "[亲亲]", "[色]", "[憧憬]", "[舔屏]", "[坏笑]",
		"[阴险]", "[笑而不语]", "[偷笑]", "[酷]", "[并不简单]", "[思考]", "[疑问]", "[费解]", "[晕]", "[衰]",
		"[骷髅]", "[嘘]", "[闭嘴]", "[傻眼]", "[吃惊]", "[吐]", "[感冒]", "[生病]", "[拜拜]", "[鄙视]", "[白眼]",
		"[左哼哼]", "[右哼哼]", "[抓狂]", "[怒骂]", "[打脸]", "[顶]", "[互粉]", "[钱]", "[哈欠]", "[困]", "[睡]",
		"[吃瓜]", "[doge]", "[二哈]", "[喵喵]", "[赞]", "[good]", "[ok]", "[耶]", "[握手]", "[作揖]", "[来]", "[拳头]",
		"[心]", "[伤心]", "[鲜花]", "[男孩儿]", "[女孩儿]", "[熊猫]", "[兔子]", "[猪头]", "[草泥马]", "[奥特曼]", "[太阳]",
		"[月亮]", "[浮云]", "[下雨]", "[沙尘暴]", "[微风]", "[围观]", "[飞机]", "[照相机]", "[话筒]", "[蜡烛]", "[音乐]",
		"[喜]", "[给力]", "[威武]", "[干杯]", "[蛋糕]", "[礼物]", "[钟]", "[肥皂]", "[绿丝带]", "[围脖]", "[浪]", "[羞嗒嗒]",
		"[好爱哦]", "[偷乐]", "[赞啊]", "[笑哈哈]", "[好喜欢]", "[求关注]", "[胖丁微笑]", "[弱]", "[NO]", "[haha]", "[加油]",
		"[佩奇]", "[大侦探皮卡丘微笑]", "[圣诞老人]", "[紫金草]", "[文明遛狗]", "[神马]", "[马到成功]", "[炸鸡啤酒]", "[最右]",
		"[织]", "[五仁月饼]", "[给你小心心]", "[吃狗粮]", "[弗莱见钱眼开]", "[点亮橙色]", "[超新星全运会]", "[看涨]", "[看跌]",
		"[带着微博去旅行]", "[星星]", "[半星]", "[空星]", "[蕾伊]", "[凯洛伦]", "[BB8]", "[冲锋队员]", "[达斯维达]", "[C3PO]",
		"[丘巴卡]", "[R2D2]", "[哆啦A梦花心]", "[哆啦A梦害怕]", "[哆啦A梦吃惊]", "[哆啦A梦汗]", "[哆啦A梦微笑]", "[伴我同行]",
		"[静香微笑]", "[大雄微笑]", "[胖虎微笑]", "[小夫微笑]", "[哆啦A梦笑]", "[哆啦A梦无奈]", "[哆啦A梦美味]", "[哆啦A梦开心]",
		"[哆啦A梦亲亲]", "[小黄人微笑]", "[小黄人剪刀手]", "[小黄人不屑]", "[小黄人高兴]", "[小黄人惊讶]", "[小黄人委屈]",
		"[小黄人坏笑]", "[小黄人白眼]", "[小黄人无奈]", "[小黄人得意]",
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
		return nil, errors.Wrap(err, "weiboclock NewClock PCLogin error")
	}
	code, err := weibo.AuthCode()
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock NewClock AuthCode error")
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

// OclockText 返回报时文本内容
func (c *Clock) OclockText() string {
	rand.Seed(time.Now().Unix())
	mood := Moods[rand.Intn(len(Moods))]
	oclock := time.Now().Hour()
	if oclock > 12 {
		oclock = oclock - 12
	}
	words := strings.Repeat(Words[rand.Intn(len(Words))], oclock)
	return fmt.Sprintf("%d点啦~  %s %s \nhttp://%s", oclock, mood, words, c.securityDomain)
}

// Toll 发送整点报时微博
func (c *Clock) Toll() error {
	if err := c.UpdateToken(); err != nil {
		return errors.Wrap(err, "weiboclock Toll UpdateToken error")
	}
	text := c.OclockText()
	if err := c.weibo.StatusesShare(c.token.AccessToken, text, nil); err != nil {
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
		code, err := c.weibo.AuthCode()
		if err != nil {
			return errors.Wrap(err, "weiboclock UpdateToken AuthCode error")
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
