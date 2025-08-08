package proxy

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	certCache = make(map[string]*tls.Certificate)
	caCert    *x509.Certificate
	caKey     crypto.PrivateKey

	privKey   *rsa.PrivateKey //私钥,不用每次生成
	transport *http.Transport //全局复用
)

func init() {
	transport = &http.Transport{
		MaxIdleConns:          100,
		MaxConnsPerHost:       10,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	privKey, _ = rsa.GenerateKey(rand.Reader, 2048)
}

func (p *Proxy) genCert(hostname string) (*tls.Certificate, error) {
	hostname = strings.Split(hostname, ":")[0]
	// 检查缓存
	if cert, ok := certCache[hostname]; ok {
		return cert, nil
	}

	// 生成私钥
	//privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to generate private key: %v", err)
	//}

	// 创建证书模板
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"ProxyDev Technology Co., Ltd"},
			CommonName:   hostname,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{hostname},
	}

	// 使用CA签名证书
	certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &privKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	// 创建tls.Certificate
	cert := &tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privKey,
	}

	// 缓存证书
	certCache[hostname] = cert
	return cert, nil
}

// MITM中间人代理
func (p *Proxy) handleHTTPSWithMITM(w http.ResponseWriter, r *http.Request) {
	p.logger.Printf("proxy https mitm: %s\n", r.URL)

	// 1. 劫持连接
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
	p.AddConnSess(clientConn)
	//clientConn.SetDeadline(time.Now().Add(3 * time.Second)) //连接自动断开

	// 必须立即发送200响应
	if _, err := clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		p.logger.Printf("Failed to send 200 response: %v\n", err)
		return
	}

	// 2. 生成目标主机证书
	cert, err := p.genCert(r.Host)
	if err != nil {
		p.logger.Printf("Failed to generate cert for %s: %v\n", r.Host, err)
		return
	}

	// 3. 与客户端进行TLS握手
	tlsConn := tls.Server(clientConn, &tls.Config{
		Certificates: []tls.Certificate{*cert},
	})
	defer tlsConn.Close()
	p.AddConnSess(tlsConn)

	if err := tlsConn.Handshake(); err != nil {
		p.logger.Printf("TLS handshake with client failed: %v\n", err)
		return
	}

	// 4. 现在可以读取和修改明文HTTP请求
	clientReader := bufio.NewReader(tlsConn)
	for {
		req, err := http.ReadRequest(clientReader)
		if err != nil {
			if err != io.EOF {
				p.logger.Printf("Read request error: %v\n", err)
			}
			return
		}

		req.URL.Scheme = "https"
		req.URL.Host = r.Host
		req.Host = r.Host

		// 请求拦截器
		if p.reqHandler != nil {
			p.reqHandler(w, req)
		}

		// 5. 转发请求到目标服务器并修改响应
		//transport := &http.Transport{
		//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//}
		resp, err := transport.RoundTrip(req)
		if err != nil {
			p.logger.Printf("Failed to round trip: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// 6. 修改响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			p.logger.Printf("Failed to read body: %v\n", err)
			return
		}

		// 相应拦截器
		modifiedBody := body
		if p.respHandler != nil {
			modifiedBody = p.respHandler(resp, body)
		}

		// 7. 写回修改后的响应
		resp.Body = io.NopCloser(bytes.NewReader(modifiedBody))
		resp.ContentLength = int64(len(modifiedBody))
		resp.Header.Set("Content-Length", fmt.Sprint(len(modifiedBody)))

		if err := resp.Write(tlsConn); err != nil {
			p.logger.Printf("Failed to write response: %v\n", err)
		}
	}
}
