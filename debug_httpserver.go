package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// 微博报时
func tollHandler(w http.ResponseWriter, r *http.Request) {
	homeURL := doToll(true)
	fmt.Fprintln(w, fmt.Sprintf(`<p>toll done. <a href="%s">%s</a></p>`, homeURL, homeURL))
}

// 首页
func indexHandler(w http.ResponseWriter, r *http.Request) {
	index := `<a href="/toll" target="blank">toll</a>`
	fmt.Fprintln(w, index)
}

// runDebugHTTPServer 运行debug http服务，提供HTTP接口调用函数
func runDebugHTTPServer(addr, basicAuthUsername, basicAuthPasswd string) {
	if addr == "" {
		addr = ":1733" // 写这里时的当前时间作为默认地址端口
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicAuth(indexHandler, basicAuthUsername, basicAuthPasswd, "debug"))
	mux.HandleFunc("/toll", basicAuth(tollHandler, basicAuthUsername, basicAuthPasswd, "weiboclock"))

	log.Println("[INFO] debug http server is running at", addr2URL(addr))
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Println("[ERROR] runDebugHTTPServer ListenAndServe error:", err)
	}
}

// basicAuth 为handlefunc注册基础认证
// curl --user name:password http://www.example.com
func basicAuth(handler http.HandlerFunc, basicAuthUsername, basicAuthPasswd, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(basicAuthUsername)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(basicAuthPasswd)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("You are Unauthorized to access the application.\n"))
			return
		}
		handler(w, r)
	}
}

// realip 获取ip地址
func realip() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("[WARN] get ip error:" + err.Error())
		return ip
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}

// addr转url
func addr2URL(addr string) string {
	s := strings.Split(addr, ":")
	if len(s) == 2 {
		port := s[1]
		return fmt.Sprintf("http://%s:%s", realip(), port)
	}
	return addr
}
