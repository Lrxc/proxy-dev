package mod

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
	"log"
	"math/big"
	"net/http"
	"net/http/httputil"
	resources "proxy-dev/assets"
	"strings"
	"testing"
	"time"
)

var (
	caCert    *x509.Certificate
	caKey     crypto.PrivateKey
	certCache = make(map[string]*tls.Certificate)
)

func initCa() {
	certFile, _ := resources.ReadByte("server.crt")
	keyFile, _ := resources.ReadByte("server.key")
	goproxyCa, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	caCert = goproxyCa.Leaf
	caKey = goproxyCa.PrivateKey
}

func TestMod(t *testing.T) {
	initCa()

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req", r.Host)

		if r.Method == http.MethodConnect {
			handleHTTPS(w, r)
		} else {
			handleHTTP(w, r)
		}
	}

	// 启动代理服务器
	server := &http.Server{
		Addr:    ":10086",
		Handler: http.HandlerFunc(handler),
	}

	log.Println("MITM Proxy started on :10086")
	log.Println("Please install ca-cert.pem as trusted root CA")
	log.Fatal(server.ListenAndServe())
}

func genCert(hostname string) (*tls.Certificate, error) {
	hostname = strings.Split(hostname, ":")[0]

	// 检查缓存
	if cert, ok := certCache[hostname]; ok {
		return cert, nil
	}

	// 生成私钥
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// 创建证书模板
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"MITM Proxy"},
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

func modifyContent(body []byte) []byte {
	// 这里实现你的内容修改逻辑
	// 示例: 替换所有"Google"为"MITM Proxy"
	modified := bytes.ReplaceAll(body, []byte("Google"), []byte("MITM Proxy"))
	return modified
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	// 创建Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 创建反向代理
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = r.Host
			req.Host = r.Host
		},
		Transport: transport,
		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode == http.StatusOK {
				contentType := resp.Header.Get("Content-Type")
				if strings.Contains(contentType, "text/html") {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return err
					}
					resp.Body.Close()

					modifiedBody := modifyContent(body)
					resp.Body = io.NopCloser(bytes.NewReader(modifiedBody))
					resp.ContentLength = int64(len(modifiedBody))
					resp.Header.Set("Content-Length", fmt.Sprint(len(modifiedBody)))
				}
			}
			return nil
		},
	}

	proxy.ServeHTTP(w, r)
}

func handleHTTPS(w http.ResponseWriter, r *http.Request) {
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

	// 必须立即发送200响应
	if _, err := clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		log.Printf("Failed to send 200 response: %v", err)
		return
	}

	// 2. 生成目标主机证书
	cert, err := genCert(r.Host)
	if err != nil {
		log.Printf("Failed to generate cert for %s: %v", r.Host, err)
		return
	}

	// 3. 与客户端进行TLS握手
	tlsConn := tls.Server(clientConn, &tls.Config{
		Certificates: []tls.Certificate{*cert},
	})
	defer tlsConn.Close()

	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake with client failed: %v", err)
		return
	}

	// 4. 现在可以读取和修改明文HTTP请求
	req, err := http.ReadRequest(bufio.NewReader(tlsConn))
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		return
	}

	originalHost := r.Host
	//重定向
	if r.Host == "www.baidu.com:443" {
		originalHost = "www.bing.com"
	}
	req.URL.Scheme = "https"
	req.URL.Host = originalHost
	req.Host = originalHost

	// 5. 转发请求到目标服务器并修改响应
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Printf("Failed to round trip: %v", err)
		return
	}
	defer resp.Body.Close()

	// 6. 修改响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		return
	}

	modifiedBody := modifyContent(body)

	// 7. 写回修改后的响应
	resp.Body = io.NopCloser(bytes.NewReader(modifiedBody))
	resp.ContentLength = int64(len(modifiedBody))
	resp.Header.Set("Content-Length", fmt.Sprint(len(modifiedBody)))

	if err := resp.Write(tlsConn); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
