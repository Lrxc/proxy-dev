package proxy

import (
	"io"
	"net"
	"net/http"
	"proxy-dev/internal/config"
	"strings"
)

func (p *Proxy) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	if !config.Conf.System.Https {
		p.handleHTTPSDirect(w, r)
		return
	}

	if p.reqHandler != nil {
		request := r.Clone(r.Context())

		//分流判断,是否需要代理
		handler := p.shuntHandler(w, request)
		if handler {
			p.handleHTTPSWithMITM(w, r)
			return
		}
	}

	p.handleHTTPSDirect(w, r)
}

func (p *Proxy) handleHTTPSDirect(w http.ResponseWriter, r *http.Request) {
	p.logger.Printf("proxy https direct: %s\n", r.URL)

	// 获取目标主机和端口
	hostPort := r.Host
	if !strings.Contains(hostPort, ":") {
		hostPort = net.JoinHostPort(hostPort, "443")
	}

	// 与目标服务器建立连接
	destConn, err := net.Dial("tcp", hostPort)
	if err != nil {
		p.logger.Printf("连接目标服务器失败: %v\n", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer destConn.Close()
	p.AddConnSess(destConn)

	// 劫持客户端连接
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		p.logger.Printf("劫持连接失败: %v\n", err)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()
	p.AddConnSess(clientConn)

	// 告诉客户端连接已建立
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// 双向转发数据
	go p.transfer(destConn, clientConn)
	p.transfer(clientConn, destConn)
}

func (p *Proxy) transfer(destination io.Writer, source io.Reader) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Printf("转发数据时发生错误: %v\n", r)
		}
	}()

	if _, err := io.Copy(destination, source); err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			p.logger.Printf("转发数据错误: %v\n", err)
		}
	}
}
