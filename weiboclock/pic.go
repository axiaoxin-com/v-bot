// Package pic 微博图片参数
// 仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M
// pic_path为空不发送图片，为default时使用内置assets/weibo中的图片，配置中如果指定了weiboclock.pic_path，则使用指定目录中以整点数命名的png图片
// 内置图片的处理逻辑设想：
// 每次都从doutula上获取一张随机图片融入到当前内置图的表盘中央合成一张新的图，获取失败则用一个默认的图合成

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
	"github.com/rakyll/statik/fs"
)

// PicReader 返回图片的io.Reader对象
// path 为空不返回图片， default返回默认内置图片，其他返回指定路径下的%d.png命名的图片
func PicReader(path string, hour int) (io.Reader, error) {
	log.Printf("[DEBUG] PicReader path=%s hour=%d", path, hour)
	switch path {
	case "":
		return nil, nil
	case "default":
		// 使用默认内置图片
		statikFS, err := fs.New()
		if err != nil {
			return nil, errors.Wrap(err, "cuitclock PicReader New statikFS error")
		}
		filename := fmt.Sprintf("/weibo/%d.png", hour)
		f, err := statikFS.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock PicReader statikFS.Open error")
		}
		return f, nil
	default:
		// 设置了pic_path，使用自定义图片
		filename := filepath.Join(path, fmt.Sprintf("%d.png", hour))
		f, err := os.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock PicReader os.Open error")
		}
		return f, nil
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
