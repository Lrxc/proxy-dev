package proxy

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func (p *Proxy) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 从请求头中获取目标地址
	targetURL := r.Header.Get("X-Target-URL")
	if targetURL == "" {
		// 如果没有指定目标地址，默认使用请求的Host
		targetURL = "wss://" + r.Host
	}

	// 解析目标URL
	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	// 建立到目标服务器的连接
	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", target.Host)
	if err != nil {
		http.Error(w, "Error dialing target server", http.StatusBadGateway)
		return
	}
	defer conn.Close()

	// 劫持客户端连接
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Error hijacking connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// 修改请求目标
	r.URL.Scheme = target.Scheme
	r.URL.Host = target.Host
	r.Host = target.Host

	// 移除X-Target-URL头
	r.Header.Del("X-Target-URL")

	// 发送请求到目标服务器
	err = r.Write(conn)
	if err != nil {
		p.logger.Printf("Error writing request to target: %v", err)
		return
	}

	// 双向转发数据
	go func() {
		_, err := copyData(conn, clientConn)
		if err != nil {
			p.logger.Printf("Error copying from target to client: %v", err)
		}
	}()

	_, err = copyData(clientConn, conn)
	if err != nil {
		p.logger.Printf("Error copying from client to target: %v", err)
	}
}

// 复制数据
func copyData(dst net.Conn, src net.Conn) (int64, error) {
	return io.Copy(dst, src)
}
