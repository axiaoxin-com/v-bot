package main

import (
	"cuitclock/config"
	"log"

	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	log.Println("[INFO] cuitclock inited config.")
	// 开启debug则启动一个http服务，通过请求即可立即使用测试账号发微博
	if viper.GetBool("server.debug") {
		// debug 模式启动debug http server
		go runDebugHTTPServer(viper.GetString("server.debug_addr"), viper.GetString("server.basic_auth_username"), viper.GetString("server.basic_auth_passwd"))
	}
	runCronServer()
}
