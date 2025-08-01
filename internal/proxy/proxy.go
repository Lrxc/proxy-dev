package proxy

import (
	"crypto/tls"
	"net/http"
	"strings"
)

type (
	ReqHandler  func(w http.ResponseWriter, r *http.Request) bool
	RespHandler func(w *http.Response, body []byte) []byte
)

// 实现了http.Handler接口
type Proxy struct {
	shuntHandler ReqHandler
	reqHandler   ReqHandler
	respHandler  RespHandler

	logger Logger
}

func NewProxy() *Proxy {
	proxy := &Proxy{}

	proxy.logger = DefaultLog{}
	return proxy
}

// ServeHTTP 实现了http.Handler接口
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.logger.Printf("proxy req %s\n", r.URL)

	switch {
	case r.Method == http.MethodConnect:
		p.handleHTTPS(w, r)
	case p.IsWebSocketUpgrade(r):
		//p.handleWebSocket(w, r)
		p.handleHTTPSDirect(w, r)
	default:
		p.handleHTTP(w, r)
	}
}

func (p *Proxy) IsWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Connection"), "upgrade") &&
		strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

// 设置自定义证书
func (p *Proxy) LoadCert(f func() tls.Certificate) {
	cert := f()
	caCert = cert.Leaf
	caKey = cert.PrivateKey
}

// 设置log
func (p *Proxy) Logger(logger Logger) {
	p.logger = logger
}

// 请求拦截器
func (p *Proxy) ShuntHandler(handler ReqHandler) {
	p.shuntHandler = handler
}

// 请求拦截器
func (p *Proxy) ReqHandler(handler ReqHandler) {
	p.reqHandler = handler
}

// 相应拦截器
func (p *Proxy) RespHandler(handler RespHandler) {
	p.respHandler = handler
}
