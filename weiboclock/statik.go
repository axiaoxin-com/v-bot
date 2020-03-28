package weiboclock

import (
	"log"
	"net/http"

	// 导入statik生成的代码
	_ "cuitclock/statik"

	"github.com/rakyll/statik/fs"
)

// StatikFS statik obj
var StatikFS http.FileSystem

func init() {
	var err error
	StatikFS, err = fs.New()
	if err != nil {
		log.Fatalln("weiboclock init StatikFS error:", err)
	}
}
