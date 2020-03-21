// Package pic 微博图片参数
// 仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M
// pic_path为空不发送图片，为default时使用内置assets/weibo中的图片，配置中如果指定了weiboclock.pic_path，则使用指定目录中以整点数命名的png图片
// 内置图片的处理逻辑设想：
// 每次都从doutula上获取一张随机图片融入到当前内置图的表盘中央合成一张新的图，获取失败则用一个默认的图合成

package weiboclock

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

// PicReader 返回图片的io.Reader对象
// path 为空不返回图片， default返回默认内置图片，其他返回指定路径下的%d.png命名的图片
func PicReader(path string, hour int) (io.Reader, error) {
	log.Printf("[DEBUG] PicReader path=%s hour=%d\n", path, hour)
	switch path {
	case "":
		return nil, nil
	case "default":
		// 使用内置表盘图片生成上传的图片
		statikFS, err := fs.New()
		if err != nil {
			return nil, errors.Wrap(err, "cuitclock PicReader New statikFS error")
		}
		filename := fmt.Sprintf("/weibo/%d.png", hour)
		f, err := statikFS.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock PicReader statikFS.Open error")
		}
		// 获取doutula表情包
		picURLs, err := DoutulaSearch(strconv.Itoa(hour), 1)
		if err != nil {
			// 获取失败则使用默认图片
			log.Println("[ERROR] weiboclock PicReader DoutulaSearch error:" + err.Error())
			return f, nil
		}
		f1, format, err := PickOnePicFromURLs(picURLs)
		if err != nil {
			log.Println("[ERROR] weiboclock PicReader PickOnePicFromURLs error:" + err.Error())
			// 获取失败则使用默认图片
			icon, err := os.Open("/weibo/icon.jpg")
			if err != nil {
				// 默认图片失败则使用最原始的图片
				return f, nil
			}
			f1 = icon
			format = "jpg"
		}
		defer f1.Close()
		// 将获取的图片融合到表盘中央
		mf, err := MergeClockPic(f, f1, format)
		if err != nil {
			// 融合失败则使用默认图片
			log.Println("[ERROR] weiboclock PicReader MergeClockPic error:" + err.Error())
			return f, nil
		}
		return mf, nil
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
	allowFormats := []string{"jpg", "jpeg", "png", "gif"}
	// 正则匹配出图片url
	picSrc, err := regexp.Compile(`<img referrerpolicy="no-referrer".+ data-original="(.+?)"`)
	for _, matched := range picSrc.FindAllStringSubmatch(string(body), -1) {
		if len(matched) == 2 {
			picURL := matched[1]
			// 过滤出指定格式的图片
			picURLSplited := strings.Split(picURL, ".")
			format := picURLSplited[len(picURLSplited)-1]
			for _, allowFormat := range allowFormats {
				if strings.ToLower(format) == allowFormat {
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

// MergeClockPic 合并表盘和获取的图片
// 参考文章：https://blog.golang.org/image-draw
// https://golang.org/doc/progs/image_draw.go
func MergeClockPic(clock, pic io.Reader, format string) (*bytes.Buffer, error) {
	var background image.Image
	var front image.Image
	var err error

	background, err = png.Decode(clock)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode clock error")
	}

	switch format {
	case "png":
		front, err = png.Decode(pic)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode pic as png error")
		}
	case "jpg", "jpeg":
		front, err = jpeg.Decode(pic)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode pic as jpeg error")
		}
	case "gif":
		front, err = gif.Decode(pic)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode pic as gif error")
		}
	}

	frontWidth := 300
	frontHeight := frontWidth
	front = resize.Resize(uint(frontWidth), uint(frontHeight), front, resize.Lanczos3)

	ftBounds := front.Bounds()
	bgBounds := background.Bounds()
	// front 放表盘中央
	ftOffsetWidth := bgBounds.Size().X/2 - int(frontWidth)/2
	ftOffsetHeight := bgBounds.Size().Y/2 - int(frontHeight)/2 + 40 // 不+40不能在中心位置
	frontOffset := image.Pt(ftOffsetWidth, ftOffsetHeight)

	// front 画成圆形
	p := image.Pt(frontWidth/2, frontHeight/2)
	r := frontWidth / 2
	circle := &Circle{p, r}

	img := image.NewRGBA(bgBounds)
	draw.Draw(img, bgBounds, background, bgBounds.Min, draw.Src)
	draw.DrawMask(img, ftBounds.Add(frontOffset), front, ftBounds.Min, circle, ftBounds.Min, draw.Over)

	imgBuf := new(bytes.Buffer)
	err = png.Encode(imgBuf, img)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock MergeClockPic png.Encode error")
	}
	return imgBuf, nil
}

// Circle 圆形Mask
type Circle struct {
	p image.Point
	r int
}

// ColorModel godoc
func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

// Bounds godoc
func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

// At godoc
func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
