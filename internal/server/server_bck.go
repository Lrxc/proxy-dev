package server

//
//import (
//	"crypto/tls"
//	"crypto/x509"
//	"fmt"
//	"github.com/elazarl/goproxy"
//	log "github.com/sirupsen/logrus"
//	"net/http"
//	resources "proxy-dev/assets"
//	"proxy-dev/internal/config"
//)
//
//var server *http.Server
//
//func StartServer(https bool) {
//	proxy := goproxy.NewProxyHttpServer()
//	//proxy.Verbose = true
//	proxy.AllowHTTP2 = true
//	proxy.Logger = SilentLog{}
//
//	if https {
//		// 设置CA证书
//		if err := SetCA(); err != nil {
//			log.Fatalf("https certificate err: %v", err)
//		}
//
//		// 处理HTTPS连接
//		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
//	}
//
//	// 请求
//	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
//		fmt.Println(req.URL.String())
//		return req, nil
//	})
//
//	addr := fmt.Sprintf("%s:%d", config.Conf.Proxy.Host, config.Conf.Proxy.Port)
//	log.Info("server listen: ", addr)
//
//	server = &http.Server{Addr: addr, Handler: proxy}
//	server.ListenAndServe()
//	//log.Fatal(http.ListenAndServe(addr, proxy))
//}
//
//func ReStartServer(b bool) error {
//	if server == nil {
//		return nil
//	}
//
//	err := StopServer()
//	go StartServer(b)
//
//	return err
//}
//
//func StopServer() error {
//	if server == nil {
//		return fmt.Errorf("服务未启动")
//	}
//
//	if err := server.Close(); err != nil {
//		return err
//	}
//
//	log.Info("server stopped")
//	return nil
//}
//
//func SetCA() error {
//	certFile, _ := resources.ReadByte("server.crt")
//	keyFile, _ := resources.ReadByte("server.key")
//
//	goproxyCa, err := tls.X509KeyPair(certFile, keyFile)
//	if err != nil {
//		return err
//	}
//	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
//		return err
//	}
//	goproxy.GoproxyCa = goproxyCa
//
//	//goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
//	//goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
//	//goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
//	//goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
//
//	return nil
//}
//
//type SilentLog struct{}
//
//func (e SilentLog) Printf(format string, v ...any) {
//	//log.Infof("goproxy: "+format, v...)
//}
