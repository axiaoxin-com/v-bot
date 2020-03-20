// Package pic 微博图片方案
// 仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M
// 根据配置文件获取图片方案
// local类型时从配置中指定的目录中获取整点数命名的png图片
// time类型时通过代码生成一张写有当前时间2006-01-02 15:00:00的图片
// online类型时按报时的数字去请求https://www.doutula.com/search?keyword=1 从其中的表情包中随机选择一张符合要求的图作为本次报时的上传图片，获取失败时使用time保底
// 没有指定方案时不使用图片

package weiboclock

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// PicReader 返回图片的io.Reader对象
func PicReader(plan, path string, oclock int) (io.Reader, error) {
	log.Printf("[DEBUG] PicReader plan=%s path=%s oclock=%d", plan, path, oclock)
	switch strings.ToLower(plan) {
	case "local":
		filename := filepath.Join(path, fmt.Sprintf("%d.png", oclock))
		f, err := os.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "weiboclock PicReader Open error")
		}
		return f, nil
	default:
		return nil, nil
	}
}
