package test

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

// https://pyer.dev/post/implement-reverse-proxy-with-golang
func TestSimple(t *testing.T) {
	target, err := url.Parse("https://www.baidu.com")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		//https需要配置
		req.Host = target.Host
	}

	http.ListenAndServe(":10086", proxy)
}
