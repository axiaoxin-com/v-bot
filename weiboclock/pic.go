// 生成微博图片参数
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
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

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

// PicReader 返回图片的io.Reader对象
// path 为空不返回图片， default返回默认内置图片，其他返回指定路径下的%d.png命名的图片
func PicReader(path string, now time.Time) (io.Reader, error) {
	hour := now.Hour()
	log.Printf("[DEBUG] PicReader path=%s hour=%d\n", path, hour)
	switch path {
	case "":
		return nil, nil
	case "default":
		// 获取表盘图片，失败直接返回
		clockPic, err := ClockPic(hour)
		if err != nil {
			return nil, err
		}
		// 获取当前时间的表盘中心位置图片
		centerPic, centerPicFormat, centerPicBgColor := CenterPic(hour)
		// 将中心图片融合到表盘中央
		mergedPic, err := MergeClockPic(now, clockPic, centerPic, centerPicFormat, centerPicBgColor)
		if err != nil {
			// 融合失败则直接返回表盘背景图片
			log.Println("[ERROR] weiboclock PicReader MergeClockPic error:", err)
			return clockPic, nil
		}
		return mergedPic, nil
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

// ClockPic 获取当前时间的表盘图片
func ClockPic(hour int) (io.ReadCloser, error) {
	filename := fmt.Sprintf("/images/clock/%d.png", hour)
	// 表盘图片
	f, err := StatikFS.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock ClockPic StatikFS.Open error")
	}
	return f, nil
}

// CenterPic 获取中心位置图片
func CenterPic(hour int) (io.ReadCloser, string, color.RGBA) {
	// 初始化声明为icon图片，后面会各种尝试覆盖该图片使用其他图片
	centerPicFormat := "jpg"
	centerPicBgColor := color.RGBA{255, 255, 255, 255} // 白色背景
	icon, err := StatikFS.Open("/images/clock/icon.jpg")
	if err != nil {
		// 获取失败后续逻辑会尝试覆盖
		log.Println("[ERROR] PicReader open icon error", err)
	}
	// 将File类型强制转为readcloser类型
	centerPic := io.ReadCloser(icon)
	defer centerPic.Close()

	// 尝试1. 使用天气图片作为中心位置图片
	if WttrInImage != nil {
		centerPic = WttrInImage
		centerPicFormat = "png"
		centerPicBgColor = color.RGBA{0, 0, 0, 255} // 黑色背景
	} else {
		// 尝试2. 天气图片尝试失败则尝试使用斗图啦表情
		log.Println("[WARN] PicReader get a nil WttrInImage, try to use doutulaPic")
		if picURLs, err := DoutulaSearch(strconv.Itoa(hour), 1); err == nil {
			// 成功搜索表情包后随机选择一个表情图片
			if doutulaPic, doutulaPicFormat, err := PickOnePicFromURLs(picURLs); err == nil {
				// 成功获取到随机表情图片后设置其为centerPic
				// 背景采用默认的白色
				centerPic = doutulaPic
				centerPicFormat = doutulaPicFormat
			} else {
				log.Println("[ERROR] PicReader PickOnePicFromURLs error", err)
			}
		} else {
			log.Println("[ERROR] PicReader DoutulaSearch error", err)
		}
	}
	return centerPic, centerPicFormat, centerPicBgColor
}

// MergeClockPic 合并表盘和获取的图片
// 参考文章：https://blog.golang.org/image-draw
// https://golang.org/doc/progs/image_draw.go
func MergeClockPic(now time.Time, clockPic, centerPic io.ReadCloser, centerPicFormat string, centerPicBgColor color.RGBA) (*bytes.Buffer, error) {
	var background image.Image
	var front image.Image
	var err error

	// 背景表盘
	background, err = png.Decode(clockPic)
	if err != nil {
		return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode clockPic error")
	}
	bgBounds := background.Bounds()

	// 中心位置图片
	switch centerPicFormat {
	case "png":
		front, err = png.Decode(centerPic)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode centerPic as png error")
		}
	case "jpg", "jpeg":
		front, err = jpeg.Decode(centerPic)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode centerPic as jpeg error")
		}
	case "gif":
		front, err = gif.Decode(centerPic)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock MergeClockPic Decode centerPic as gif error")
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
	// 设置画布背景为中心图片的背景色
	for m := 0; m < canvasWidth; m++ {
		for n := 0; n < canvasHeight; n++ {
			canvas.SetRGBA(m, n, centerPicBgColor)
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

		text := now.Format("2006-01-02")
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
