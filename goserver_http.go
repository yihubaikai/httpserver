package main

import (
	"fmt"
	"github.com/yihubaikai/gopublic"
	"net/http"
	"strings"
)

func getCurrentIP(r http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]

	}
	return ip
}

func Route_Home(w http.ResponseWriter, r *http.Request) {
	sRet := hPub.Gettime() + ":" + r.URL.String() + ":" + getCurrentIP(*r)
	fmt.Fprintln(w, sRet)
	fmt.Println(sRet)
}

func main() {
	fmt.Println("LISTENING:80")
	http.HandleFunc("/", Route_Home)
	http.ListenAndServe(":80", nil)
}
