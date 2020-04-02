package weiboclock

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

// RandFont 随机返回fonts中的一个字体
func RandFont() (*truetype.Font, error) {
	fontPaths := []string{}
	fs.Walk(StatikFS, "/fonts", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fontPaths = append(fontPaths, path)
		}
		return nil
	})
	rand.Seed(time.Now().UnixNano())
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
