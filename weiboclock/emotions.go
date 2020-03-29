package weiboclock

import (
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// TollVoices 钟楼报时声音
	TollVoices = []string{
		"biu! ", "ha! ", "mua~ ", "dang! ", "ho! ", "hei! ", "pia! ", "dong~ ",
		"he! ", "mia~ ", "ao~ ", "do~ ", "ga~ ", "bi! ", "ba~ ", "ma~ ", "miao~ ",
		"pa! ", "no~ ", "bomb! ", "yeah~ ", "ka! ", "la~ ", "da! ", "OA! ", "meow~ ",
	}

	// TextEmotions 颜文字表情
	TextEmotions = []string{
		"(°ー°〃)", "_(:з」∠)_ ", "o(*≧▽≦)ツ┏━┓", "๑乛◡乛๑ ", "(σ‘・д･)σ", "( ＿ ＿)ノ｜", "┑(￣Д ￣)┍",
		"(＃°Д°)", "(-ω- )", "(′・ω・`)", "( ^ω^)", "乀(ˉεˉ乀)", "ʕ•̫͡•ʕ•̫͡•ʔ•̫͡•ʔ", "(๑‾᷆д‾᷇๑)Fightᵎᵎ", "꒰｡⁻௰⁻｡꒱",
		"\\(╯-╰)/", "\\(▔▽▔)/", "^(oo)^", "(O^~^O)", "（╯＾╰）", "(*ﾟ∀ﾟ*)", "(శωశ)", "(*´∀`)~♥", "σ`∀´)σ",
		"(〃∀〃)", "(^_っ^)", "(｡◕∀◕｡)", "ヽ(✿ﾟ▽ﾟ)ノ", "ε٩(๑> ₃ <)۶з", "(σ′▽‵)′▽‵)σ", "σ ﾟ∀ ﾟ) ﾟ∀ﾟ)σ",
		"｡:.ﾟヽ(*´∀`)ﾉﾟ.:｡", "(✪ω✪)", "(∂ω∂)", "─=≡Σ((( つ•̀ω•́)つ", "(๑´ڡ`๑)", "(´▽`ʃ♡ƪ)", "(❛◡❛✿)", "(灬ºωº灬)",
		"(￣▽￣)/", "╰(*°▽°*)╯", "(๑•̀ㅂ•́)و✧", "( ^ω^)", "٩(｡・ω・｡)و", "( ～'ω')～", "(๑ơ ₃ ơ)♥", "(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧",
		"o(☆Ф∇Ф☆)o", "(￫ܫ￩)", "(♥д♥)", "✧◝(⁰▿⁰)◜✧", "(ᗒᗨᗕ)/", "(=´ω`=)", "(｢･ω･)｢", "(*´д`)", "Σ>―(〃°ω°〃)♡→",
		"(▰˘◡˘▰)", "ヾ(´ε`ヾ)", "(っ●ω●)っ", "◥(ฅº￦ºฅ)◤", "ヽ( ° ▽°)ノ", "(　ﾟ∀ﾟ) ﾉ♡", "✧*｡٩(ˊᗜˋ*)و✧*｡",
		"⁽⁽◟(∗ ˊωˋ ∗)◞ ⁾⁾", "ヾ(´︶`*)ﾉ♬", "ヾ(*´∀ ˋ*)ﾉ", "(๑•̀ω•́)ノ", "ヾ (o ° ω ° O ) ノ゙ ", "╮(╯_╰)╭", "(๑•́ ₃ •̀๑)",
		"(´･_･`)", "(ㆆᴗㆆ)", "┐(´д`)┌", "( ˘･з･)", "( ´•︵•` )", "(｡ŏ_ŏ)", "(◞‸◟)", "( ˘•ω•˘ )", "(눈‸눈)", "(´･ω･`)",
		"(*´艸`*)", "(〃∀〃)", "(つд⊂)", "(๑´ㅂ`๑)", "ε٩(๑> ₃ <)۶з", "(๑´ڡ`๑)", "(灬ºωº灬)", "(๑• . •๑)",
		"(๑ơ ₃ ơ)♥", "(●｀ 艸´)", ",,Ծ‸Ծ,,", "(〃ﾟдﾟ〃)", "(๑´ㅁ`)", "(๑¯∀¯๑)", "(〃´∀｀)", "(⋟﹏⋞)", "(ノдT)",
		"(T_T)", "：ﾟ(｡ﾉω＼｡)ﾟ･｡", "(TдT)", "(☍﹏⁰)", "(╥﹏╥)", "｡ﾟ(ﾟ´ω`ﾟ)ﾟ｡", "இдஇ", "｡ﾟヽ(ﾟ´Д`)ﾉﾟ｡", "。･ﾟ･(つд`ﾟ)･ﾟ･",
		"・゜・(PД`q｡)・゜・", "(ﾟд⊙)", "(‘⊙д-)", "Σ(*ﾟдﾟﾉ)ﾉ", "(((ﾟДﾟ;)))", "(((ﾟдﾟ)))", "(☉д⊙)", "(|||ﾟдﾟ)",
		"(´⊙ω⊙`)", "ฅ(๑*д*๑)ฅ!!", "(゜ロ゜)", "(✘﹏✘ა)", "(✘Д✘๑ )", "(╬☉д⊙)", "(／‵Д′)／~ ╧╧", "(╯‵□′)╯︵┴─┴",
		"(◓Д◒)✄╰⋃╯", "(ﾒﾟДﾟ)ﾒ", "(`へ´≠)", "(#ﾟ⊿`)凸", "(╬▼дﾟ)", "(ᗒᗣᗕ)՞", "( ิ◕㉨◕ ิ)", "(❍ᴥ❍ʋ)", "(◕ܫ◕)", "(ΦωΦ)",
		"ก็ʕ•͡ᴥ•ʔ ก้", "(=´ω`=)", "(⁰⊖⁰)", "(=´ᴥ`)", "ฅ●ω●ฅ", "( ° ͜ʖ͡°)╭∩╮", "(⌐▀͡ ̯ʖ▀)", "(･ิω･ิ)", "ʕ•̀ω•́ʔ✧", "٩(♡ε♡ )۶", "٩(๑´3｀๑)۶",
		"(๑•̀ㅁ•́๑)✧", "•̀.̫•́✧", "⁄(⁄ ⁄•⁄ω⁄•⁄ ⁄)⁄", "⁽⁽٩(๑˃̶͈̀ ᗨ ˂̶͈́)۶⁾⁾", "( •̀ᄇ• ́)ﻭ✧", "(▭-▭)✧", "(▭-▭)✧", "ଘ(⊃ˊᵕˋ)⊃━☆’･*:.｡",
		"(⸝⸝⸝ᵒ̴̶̥́ ⌑ ᵒ̴̶̣̥̀⸝⸝⸝)", "((̵̵́ ̆͒͟˚̩̭ ̆͒)̵̵̀)ﾞ", "o̖⸜((̵̵́ ̆͒͟˚̩̭ ̆͒)̵̵̀)⸝o̗", "٩(ˊᗜˋ*)و", "(∩ᵒ̴̶̷̤⌔ᵒ̴̶̷̤∩)", "₍₍ ◟꒰ ‾᷅д̈ ‾᷄ ╬꒱", "ଘ(੭ˊ꒳ˋ)੭✧", "ʕ·͡·̫͖ʕ⁎̯͡⁎ʔ",
	}

	// WeiboEmotions 微博官方表情，weiboclock Run方法调用时进行初始化
	WeiboEmotions = []string{}
)

// AllEmotions 返回全部表情列表
func AllEmotions() []string {
	return append(TextEmotions, WeiboEmotions...)
}

// PickOneEmotion 随机选择一个表情
func PickOneEmotion() string {
	rand.Seed(time.Now().Unix())
	emotions := AllEmotions()
	return emotions[rand.Intn(len(emotions))]
}

// TollVoice 报时拟声
func TollVoice(count int) string {
	rand.Seed(time.Now().Unix())
	voice := TollVoices[rand.Intn(len(TollVoices))]
	return strings.Repeat(voice, count)
}

// InitWeiboEmotions 初始化微博官方表情
func (clock *WeiboClock) InitWeiboEmotions() (int, error) {
	// reset
	WeiboEmotions = []string{}
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
