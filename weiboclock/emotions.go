package weiboclock

import (
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// TollTails 标点小尾巴
	TollTails = []string{
		"!", "~", ".", "?",
	}

	// ClockEmoji 整点emoji
	ClockEmoji = []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"}

	// WeiboEmotions 微博官方表情，weiboclock Run方法调用时进行初始化
	WeiboEmotions = []string{}
)

// PickOneEmotion 随机选择一个表情
func PickOneEmotion() string {
	rand.Seed(time.Now().UnixNano())
	return WeiboEmotions[rand.Intn(len(WeiboEmotions))]
}

// TollTail 随机获取标点小尾巴~
func TollTail(count int) string {
	rand.Seed(time.Now().UnixNano())
	tail := TollTails[rand.Intn(len(TollTails))]
	return strings.Repeat(tail, count)
}

// InitEmotions 初始化表情，返回表情总数
func (clock *WeiboClock) InitEmotions() (int, error) {
	// reset
	WeiboEmotions = []string{}

	// 获取微博官方表情
	vb := clock.cronWeibo.WeiboClient()
	token := clock.cronWeibo.Token()
	language := "cnname"
	emotionType := "face"
	emotions, err := vb.Emotions(token.AccessToken, emotionType, language)
	if err != nil {
		return 0, errors.Wrap(err, "weiboclock InitWeiboEmotions Emotions error")
	}
	for _, emotion := range *emotions {
		WeiboEmotions = append(WeiboEmotions, emotion.Phrase)
	}
	return len(WeiboEmotions), nil
}
