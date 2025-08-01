package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	resources "proxy-dev/assets"
	"strings"
	"testing"
)

func TestMod(t *testing.T) {
	proxy := NewProxy()
	proxy.LoadCert(func() tls.Certificate {
		certFile, _ := resources.ReadByte("server.crt")
		keyFile, _ := resources.ReadByte("server.key")
		cert, _ := tls.X509KeyPair(certFile, keyFile)
		return cert
	})

	proxy.ReqHandler(func(w http.ResponseWriter, r *http.Request) bool {
		reqUrl := r.URL.String()
		if strings.Contains(reqUrl, "baidu") {
			newUrl := strings.ReplaceAll(reqUrl, "baidu", "so")
			newURL, _ := url.Parse(newUrl)

			//修改请求URL
			r.URL = newURL
			r.Host = newURL.Host

			fmt.Printf("redirect: %s %s\n", r.Method, r.URL.String())
			return true
		}
		return false
	})

	proxy.RespHandler(func(w *http.Response, body []byte) []byte {
		modified := bytes.ReplaceAll(body, []byte("Google"), []byte("MITM Proxy"))
		return modified
	})

	// 启动服务器
	server := &http.Server{
		Addr:    ":10086",
		Handler: proxy,
	}
	fmt.Println("反向代理服务启动在 :10086")
	server.ListenAndServe()
}
