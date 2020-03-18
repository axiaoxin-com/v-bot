package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	mathRand "math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// UserAgents ua list
var UserAgents []string = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36",
}

// Weibo 定义各种微博相关方法
type Weibo struct {
	client      *http.Client
	appkey      string
	appsecret   string
	redirecturi string
	username    string
	passwd      string
	userAgent   string
}

// MobileLoginResp 移动登录返回结构
type MobileLoginResp struct {
	Retcode int                    `json:"retcode"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
}

// PreLoginResp PC端prelogin登录返回结构
type PreLoginResp struct {
	Retcode    int    `json:"retcode"`
	Servertime int    `json:"servertime"`
	Pcid       string `json:"pcid"`
	Nonce      string `json:"nonce"`
	Pubkey     string `json:"pubkey"`
	Rsakv      string `json:"rsakv"`
	IsOpenlock int    `json:"is_openlock"`
	Showpin    int    `json:"showpin"`
	Exectime   int    `json:"exectime"`
}

// SsoLoginResp PC端ssologin登录返回结构
type SsoLoginResp struct {
	Retcode            string   `json:"retcode"`
	Ticket             string   `json:"ticket"`
	UID                string   `json:"uid"`
	Nick               string   `json:"nick"`
	CrossDomainURLList []string `json:"crossDomainUrlList"`
}

// RedirectResp 微博回调httpbin.org/get返回结构
type RedirectResp struct {
	Args map[string]string `json:"args"`
}

// TokenResp accesstoken返回结构
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UID         string `json:"uid"`
	IsRealName  string `json:"isRealName"`
}

// NewWeibo 创建Weibo实例
func NewWeibo(appkey, appsecret, username, passwd, redirecturi string) *Weibo {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return &Weibo{
		client:      client,
		appkey:      appkey,
		appsecret:   appsecret,
		redirecturi: redirecturi,
		username:    username,
		passwd:      passwd,
		userAgent:   randUA(),
	}
}

// 随机选一个ua
func randUA() string {
	mathRand.Seed(time.Now().Unix())
	return UserAgents[mathRand.Intn(len(UserAgents))]
}

// MobileLogin 移动端登录微博
func (w *Weibo) MobileLogin() error {
	data := url.Values{
		"username":     {w.username},
		"password":     {w.passwd},
		"savestate":    {"1"},
		"r":            {"https://weibo.cn/"},
		"ec":           {"0"},
		"pagerefer":    {"https://weibo.cn/pub/"},
		"entry":        {"mweibo"},
		"wentry":       {""},
		"loginfrom":    {""},
		"client_id":    {""},
		"code":         {""},
		"qq":           {""},
		"mainpageflag": {"1"},
		"hff":          {""},
		"hfp":          {""},
	}
	logingURL := "https://passport.weibo.cn/sso/login"
	req, err := http.NewRequest("POST", logingURL, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin NewRequest error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "deflate, br") // no gzip
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://passport.weibo.cn")
	req.Header.Set("Referer", "https://passport.weibo.cn/signin/login?entry=mweibo&r=https%3A%2F%2Fweibo.cn%2F&backTitle=%CE%A2%B2%A9&vt=")
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin Do error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin ReadAll error")
	}
	loginResp := &MobileLoginResp{}
	if err := json.Unmarshal(body, loginResp); err != nil {
		return errors.Wrap(err, "weibo MobileLogin Unmarshal error")
	}
	if loginResp.Retcode != 20000000 {
		return errors.New("weibo MobileLogin loginResp Retcode error:" + string(body))
	}
	return nil
}

// PCLogin 电脑web登录
func (w *Weibo) PCLogin() error {
	preloginURL := "https://login.sina.com.cn/sso/prelogin.php?"
	ssologinURL := "https://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.19)"
	/*
		pinURL := "https://login.sina.com.cn/cgi/pin.php"
	*/

	// 请求prelogin
	su := base64.StdEncoding.EncodeToString([]byte(w.username))
	req, err := http.NewRequest("GET", preloginURL, nil)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin NewRequest prelogin error")
	}
	params := url.Values{
		"entry":    {"weibo"},
		"su":       {su},
		"rsakt":    {"mod"},
		"checkpin": {"1"},
		"client":   {"ssologin.js(v1.4.19)"},
		"_":        {strconv.FormatInt(time.Now().UnixNano()/1e6, 10)},
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin Do prelogin error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin ReadAll prelogin error")
	}
	preLoginResp := &PreLoginResp{}
	if err := json.Unmarshal(body, preLoginResp); err != nil {
		return errors.Wrap(err, "weibo PCLogin Unmarshal preLoginResp error")
	}
	if preLoginResp.Retcode != 0 {
		return errors.New("weibo PCLogin preLoginResp Retcode error:" + string(body))
	}

	// 请求ssologin
	encMsg := []byte(fmt.Sprint(preLoginResp.Servertime, "\t", preLoginResp.Nonce, "\n", w.passwd))
	n, _ := new(big.Int).SetString(preLoginResp.Pubkey, 16)
	e, _ := new(big.Int).SetString("10001", 16)
	pubkey := &rsa.PublicKey{N: n, E: int(e.Int64())}
	sp, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, encMsg)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin EncryptPKCS1v15 error")
	}
	hexsp := hex.EncodeToString([]byte(sp))
	data := url.Values{
		"entry":      {"account"},
		"gateway":    {"1"},
		"from":       {""},
		"savestate":  {"30"},
		"useticket":  {"1"},
		"pagerefer":  {""},
		"vsnf":       {"1"},
		"su":         {su},
		"service":    {"account"},
		"servertime": {fmt.Sprint(preLoginResp.Servertime, randInt(1, 20))},
		"nonce":      {preLoginResp.Nonce},
		"pwencode":   {"rsa2"},
		"rsakv":      {preLoginResp.Rsakv},
		"sp":         {hexsp},
		"sr":         {"1536 * 864"},
		"encoding":   {"UTF - 8"},
		"cdult":      {"3"},
		"domain":     {"sina.com.cn"},
		"prelt":      {"95"},
		"returntype": {"TEXT"},
	}
	req, err = http.NewRequest("POST", ssologinURL, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin NewRequest ssologin error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin Do ssologin error")
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin ReadAll ssologin error")
	}
	ssoLoginResp := &SsoLoginResp{}
	if err := json.Unmarshal(body, ssoLoginResp); err != nil {
		return errors.Wrap(err, "weibo PCLogin Unmarshal ssoLoginResp error")
	}
	if ssoLoginResp.Retcode != "0" {
		return errors.New("weibo PCLogin ssoLoginResp Retcode error:" + string(body))
	}
	return w.loginSucceed(ssoLoginResp)
}

// randInt 产生指定数字范围内的随机数
func randInt(min int, max int) int {
	mathRand.Seed(time.Now().UnixNano())
	return min + mathRand.Intn(max-min)
}

func (w *Weibo) loginSucceed(resp *SsoLoginResp) error {
	// 请求login_url和home_url, 进一步验证登录是否成功
	s := strings.Split(strings.Split(resp.CrossDomainURLList[0], "ticket=")[1], "&ssosavestate=")
	loginURL := fmt.Sprintf("https://passport.weibo.com/wbsso/login?ticket=%s&ssosavestate=%s&callback=sinaSSOController.doCrossDomainCallBack&scriptId=ssoscript0&client=ssologin.js(v1.4.19)&_=%s", s[0], s[1], strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed NewRequest loginURL error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	res, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed Do loginURL error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed ReadAll loginURL error")
	}
	reg := regexp.MustCompile(`"uniqueid":"(.*?)"`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	if len(result) == 0 {
		return errors.New("weibo loginSucceed uniqueid pattern not match")
	}
	uid := result[0][1]
	homeURL := fmt.Sprintf("https://weibo.com/u/%s/home", uid)
	req, err = http.NewRequest("GET", homeURL, nil)
	req.Header.Set("User-Agent", w.userAgent)
	res, err = w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed Do homeURL error")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed ReadAll homeURL error")
	}
	if !strings.Contains(string(body), "我的首页") {
		return errors.New("weibo loginSucceed login failed")
	}
	return nil
}

// AuthCode 获取授权码
func (w *Weibo) AuthCode() (string, error) {
	authURL := "https://api.weibo.com/oauth2/authorize"
	referer := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", authURL, w.appkey, w.redirecturi)
	data := url.Values{
		"client_id":       {w.appkey},
		"response_type":   {"code"},
		"redirect_uri":    {w.redirecturi},
		"action":          {"submit"},
		"userId":          {w.username},
		"passwd":          {w.passwd},
		"isLoginSina":     {"0"},
		"from":            {""},
		"regCallback":     {""},
		"state":           {""},
		"ticket":          {""},
		"withOfficalFlag": {"0"},
	}
	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "weibo AuthCode NewRequest error")
	}
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "weibo AuthCode Do error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "weibo AuthCode ReadAll error")
	}
	redirectResp := &RedirectResp{}
	if err := json.Unmarshal(body, redirectResp); err != nil {
		return "", errors.Wrap(err, "weibo AuthCode Unmarshal error")
	}
	return redirectResp.Args["code"], nil
}

// AccessToken 获取token
func (w *Weibo) AccessToken(code string) (*TokenResp, error) {
	tokenURL := "https://api.weibo.com/oauth2/access_token"
	data := url.Values{
		"client_id":     {w.appkey},
		"client_secret": {w.appsecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {w.redirecturi},
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken ReadAll error")
	}
	tokenResp := &TokenResp{}
	if err := json.Unmarshal(body, tokenResp); err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken Unmarshal error")
	}
	return tokenResp, nil
}

// realip 获取ip地址
func realip() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return ip
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}

// StatusesShare 第三方分享一条链接到微博
// https://open.weibo.com/wiki/2/statuses/share
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// status	true	string	用户分享到微博的文本内容，必须做URLencode，内容不超过140个汉字，文本中不能包含“#话题词#”，同时文本中必须包含至少一个第三方分享到微博的网页URL，且该URL只能是该第三方（调用方）绑定域下的URL链接，绑定域在“我的应用 － 应用信息 － 基本应用信息编辑 － 安全域名”里设置。
// pic	false	binary	用户想要分享到微博的图片，仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M。上传图片时，POST方式提交请求，需要采用multipart/form-data编码方式。
// rip	false	string	开发者上报的操作用户真实IP，形如：211.156.0.1。
func (w *Weibo) StatusesShare(token, status string, pic io.Reader) error {
	apiURL := "https://api.weibo.com/2/statuses/share.json"
	ip := realip()
	var bodyBuf *bytes.Buffer
	var writer *multipart.Writer
	if pic == nil {
		data := url.Values{
			"access_token": {token},
			"status":       {status},
			"rip":          {ip},
		}
		bodyBuf = bytes.NewBufferString(data.Encode())
	} else {
		writer = multipart.NewWriter(bodyBuf)
		defer writer.Close()
		err := writer.WriteField("access_token", token)
		err = writer.WriteField("status", status)
		err = writer.WriteField("rip", ip)
		if err != nil {
			return errors.Wrap(err, "weibo StatusesShare WriteField error")
		}
		picWriter, err := writer.CreateFormFile("fieldname", "filename")
		io.Copy(picWriter, pic)
	}
	req, err := http.NewRequest("POST", apiURL, bodyBuf)
	if err != nil {
		return errors.Wrap(err, "weibo StatusesShare NewRequest error")
	}
	if pic == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo StatusesShare Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo StatusesShare ReadAll error")
	}
	type succResp struct {
		IDStr string `json:"idstr"`
	}
	sr := &succResp{}
	if err := json.Unmarshal(body, sr); err != nil {
		return errors.Wrap(err, "weibo StatusesShare Unmarshal error:"+string(body))
	}
	return nil
}

// TokenInfoResp get_token_info接口返回结构
type TokenInfoResp struct {
	UID      string `json:"uid"`
	Appkey   string `json:"appkey"`
	Scope    string `json:"scope"`
	CreateAt string `json:"create_at"`
	ExpireIn string `json:"expire_in"`
}

// TokenInfo 查询用户access_token的授权相关信息，包括授权时间，过期时间和scope权限
// https://api.weibo.com/oauth2/get_token_info
func (w *Weibo) TokenInfo(token string) (*TokenInfoResp, error) {
	apiURL := "https://api.weibo.com/oauth2/get_token_info"
	data := url.Values{
		"access_token": {token},
	}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo ReadAll error")
	}
	tokenInfoResp := &TokenInfoResp{}
	if err := json.Unmarshal(body, tokenInfoResp); err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo Unmarshal error")
	}
	return tokenInfoResp, nil
}
