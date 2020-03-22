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
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

// HourPic 根据hour值返回在线图片
func HourPic(hour int) (io.ReadCloser, string, error) {
	var f io.ReadCloser
	var err error
	var format string
	var picURLs []string

	// 获取doutula表情包
	picURLs, err = DoutulaSearch(strconv.Itoa(hour), 1)
	if err == nil {
		f, format, err = PickOnePicFromURLs(picURLs)
	}

	if err != nil {
		// 获取失败则使用默认图片
		icon, err := os.Open("/images/clock/icon.jpg")
		if err != nil {
			return nil, "", err
		}
		f = icon
		format = "jpg"
	}
	return f, format, nil
}

// PicReader 返回图片的io.Reader对象
// path 为空不返回图片， default返回默认内置图片，其他返回指定路径下的%d.png命名的图片
func PicReader(path string, hour int) (io.Reader, error) {
	log.Printf("[DEBUG] PicReader path=%s hour=%d\n", path, hour)
	switch path {
	case "":
		return nil, nil
	case "default":
		// 使用内置表盘图片生成上传的图片
		filename := fmt.Sprintf("/images/clock/%d.png", hour)
		f, err := StatikFS.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock PicReader StatikFS.Open error")
		}
		f1, format, err := HourPic(hour)
		if err != nil {
			return f, nil
		}
		defer f1.Close()
		// 将获取的图片融合到表盘中央
		mf, err := MergeClockPic(f, f1, format)
		if err != nil {
			// 融合失败则使用默认图片
			log.Println("[ERROR] weiboclock PicReader MergeClockPic error:", err)
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
	rand.Seed(AsiaShanghaiNow().Unix())
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

	// 背景表盘
	background, err = png.Decode(clock)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode clock error")
	}
	bgBounds := background.Bounds()

	// 中心表情包
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

	// 将表情包绘制在一个画布上，对画布进行圆角处理以显示完整表情包图片
	// 根据画布尺寸计算表情包的图片尺寸
	canvasWidth := 300 //画布宽度
	canvasHeight := canvasWidth
	radius := canvasWidth / 2 // 圆形半径
	frontWidth := int(math.Sqrt(math.Pow(float64(radius), 2) + math.Pow(float64(radius), 2)))
	frontHeight := frontWidth
	// 表情包resize
	front = resize.Resize(uint(frontWidth), uint(frontHeight), front, resize.Lanczos3)
	ftBounds := front.Bounds()

	// 创建画布
	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	// 设置画布背景为白色
	for m := 0; m < canvasWidth; m++ {
		for n := 0; n < canvasHeight; n++ {
			canvas.SetRGBA(m, n, color.RGBA{255, 255, 255, 255})
		}
	}
	canvasBounds := canvas.Bounds()
	// 计算表情包在画布上的offset
	frontOffsetX := (canvasWidth - frontWidth) / 2
	frontOffsetY := (canvasHeight - frontHeight) / 2
	frontOffset := image.Pt(frontOffsetX, frontOffsetY)
	// 将表情包画在画布上
	draw.Draw(canvas, ftBounds.Add(frontOffset), front, ftBounds.Min, draw.Over)

	// 画布圆角处理参数
	p := image.Pt(canvasWidth/2, canvasHeight/2)
	circle := &Circle{p, radius}

	// 创建最终图片
	img := image.NewRGBA(bgBounds)
	// 画上表盘背景
	draw.Draw(img, bgBounds, background, bgBounds.Min, draw.Src)
	// 添加画布
	canvasOffsetX := (bgBounds.Size().X - canvasWidth) / 2
	canvasOffsetY := (bgBounds.Size().Y-canvasHeight)/2 + 40 // +40才能在表盘中心
	canvasOffset := image.Pt(canvasOffsetX, canvasOffsetY)
	draw.DrawMask(img, canvasBounds.Add(canvasOffset), canvas, canvasBounds.Min, circle, canvasBounds.Min, draw.Over)

	// 图片底部加上当前日期
	rFont, err := RandFont()
	if err == nil {
		fc := freetype.NewContext()
		fc.SetFont(rFont) // 字体
		fc.SetDPI(72)     // 分辨率
		fontSize := 50.0
		fc.SetFontSize(fontSize)                                   //字号
		fc.SetClip(img.Bounds())                                   // 区域
		fc.SetDst(img)                                             // 目标图片
		fc.SetSrc(image.NewUniform(color.RGBA{52, 152, 219, 255})) // 字体颜色

		text := AsiaShanghaiNow().Format("2006-01-02")
		pt := freetype.Pt((bgBounds.Size().X-len(text)*25)/2, bgBounds.Size().Y-50)
		_, err := fc.DrawString(text, pt)
		if err != nil {
			log.Println("[ERROR] weiboclock MergeClockPic DrawString error:", err)
		}
	} else {
		// 字体获取失败则不加
		log.Println("[ERROR] weiboclock MergeClockPic RandFont error:", err)
	}

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

// RandFont 随机返回fonts中的一个字体
func RandFont() (*truetype.Font, error) {
	fontPaths := []string{}
	fs.Walk(StatikFS, "/fonts", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fontPaths = append(fontPaths, path)
		}
		return nil
	})
	rand.Seed(AsiaShanghaiNow().Unix())
	fontPath := fontPaths[rand.Intn(len(fontPaths))]
	log.Println("[DEBUG] weiboclock RandFont use font", fontPath)
	fontFile, err := StatikFS.Open(fontPath)
	fontBytes, err := ioutil.ReadAll(fontFile)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock RandFont StatikFS.Open error")
	}
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock RandFont ReadFile error")
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock RandFont ParseFont error")
	}
	return f, nil
}
