// Package pic 微博图片方案
// 仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M
// 根据配置文件获取图片方案
// local类型时从配置中指定的目录中获取整点数命名的png图片
// time类型时通过代码生成一张写有当前时间2006-01-02 15:00:00的图片
// online类型时按报时的数字去请求https://www.doutula.com/search?type=photo&more=1&keyword=1&page=1从其中的表情包中随机选择一张符合要求的图作为本次报时的上传图片，获取失败时使用time保底
// 没有指定方案时不使用图片

package weiboclock

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// PicReader 返回图片的io.Reader对象
func PicReader(plan, path string, hour int) (io.Reader, error) {
	log.Printf("[DEBUG] PicReader plan=%s path=%s hour=%d", plan, path, hour)
	switch strings.ToLower(plan) {
	case "local":
		filename := filepath.Join(path, fmt.Sprintf("%d.png", hour))
		f, err := os.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock PicReader Open error")
		}
		return f, nil
	default:
		return nil, nil
	}
}

// DoutulaSearch 从doutula根据关键字搜索图片，返回图片链接列表
func DoutulaSearch(keyword string, page int) ([]string, error) {
	searchURL := fmt.Sprintf("https://www.doutula.com/search?keyword=%s&type=photo&more=1&page=%d", keyword, page)
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock DoutulaSearch Get error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("weiboclock DoutulaSearch resp.Status=" + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock DoutulaSearch ReadAll error")
	}
	picURLs := []string{}
	allowSuffixs := []string{"jpg", "jpeg", "png", "gif"}
	// 正则匹配出图片url
	picSrc, err := regexp.Compile(`<img referrerpolicy="no-referrer".+ data-original="(.+?)"`)
	for _, matched := range picSrc.FindAllStringSubmatch(string(body), -1) {
		if len(matched) == 2 {
			picURL := matched[1]
			// 过滤出指定格式的图片
			picURLSplited := strings.Split(picURL, ".")
			suffix := picURLSplited[len(picURLSplited)-1]
			for _, allowSuffix := range allowSuffixs {
				if strings.ToLower(suffix) == allowSuffix {
					picURLs = append(picURLs, picURL)
				}
			}
		}
	}
	if len(picURLs) == 0 {
		return nil, errors.New("weiboclock DoutulaSearch get 0 picURLs")
	}
	return picURLs, nil
}
