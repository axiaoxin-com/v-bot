package weiboclock

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

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

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock DoutulaSearch NewDocumentFromReader error")
	}

	picURLs := []string{}
	allowFormats := []string{"jpg", "jpeg", "png", "gif"}
	dom.Find(".random_picture img[data-original]").Each(func(i int, s *goquery.Selection) {
		// 过滤出指定格式的图片
		picURL, exists := s.Attr("data-original")
		if exists {
			picURLSplited := strings.Split(picURL, ".")
			format := picURLSplited[len(picURLSplited)-1]
			for _, allowFormat := range allowFormats {
				if strings.ToLower(format) == allowFormat {
					picURLs = append(picURLs, picURL)
				}
			}
		}
	})
	if len(picURLs) == 0 {
		return nil, errors.New("weiboclock DoutulaSearch get 0 picURLs")
	}
	return picURLs, nil
}

// PickOnePicFromURLs 从给定的图片url中随机获取一张图片
func PickOnePicFromURLs(picURLs []string) (io.ReadCloser, string, error) {
	rand.Seed(time.Now().Unix())
	picURL := picURLs[rand.Intn(len(picURLs))]
	log.Println("[DEBUG] weiboclock PickOnePicFromURLs picURL:", picURL)
	picURLSplited := strings.Split(picURL, ".")
	format := picURLSplited[len(picURLSplited)-1]
	resp, err := http.Get(picURL)
	if err != nil {
		return nil, format, errors.Wrap(err, "weiboclock PickOnePicFromURLs Get error")
	}
	if resp.StatusCode != 200 {
		return nil, format, errors.Wrap(err, "weiboclock PickOnePicFromURLs StatusCode error:"+resp.Status)
	}
	return resp.Body, format, nil
}
