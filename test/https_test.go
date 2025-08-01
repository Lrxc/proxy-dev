package test

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	server := &http.Server{
		Addr: ":10086",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL.String(), "...start")

			// 处理HTTPS CONNECT请求
			if r.Method == "CONNECT" {
				handleHTTPS(w, r)
			} else {
				handleHTTP(w, r)
			}
			fmt.Println(r.URL.String(), "...end")
		}),
	}

	log.Println("Starting universal proxy server on :10086")
	log.Fatal(server.ListenAndServe())
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	// 创建动态目标代理
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			//req mod
			fmt.Println("req", req.URL)
		},
		ModifyResponse: func(resp *http.Response) error {
			//resp mod
			return nil
		},
	}

	//target, err := url.Parse(r.URL.String())
	//proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

// 处理HTTPS隧道连接
func handleHTTPS(w http.ResponseWriter, r *http.Request) {
	// 建立与目标服务器的连接
	destConn, err := net.DialTimeout("tcp", r.Host, 3*time.Second)
	destConn.SetDeadline(time.Now().Add(3 * time.Second)) //连接自动断开
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	// 响应客户端连接已建立
	w.WriteHeader(http.StatusOK)

	// 获取底层连接
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	// 双向数据转发
	go transfer(destConn, clientConn)
	transfer(clientConn, destConn)
}

// 数据转发
func transfer(destination io.Writer, source io.Reader) {
	defer destination.(io.Closer).Close()
	io.Copy(destination, source)
}
