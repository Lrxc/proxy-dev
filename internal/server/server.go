package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"proxy-dev/internal/config"
	"proxy-dev/internal/proxy"
	"proxy-dev/internal/system"
)

var (
	server *http.Server
	lproxy *proxy.Proxy
)

func StartServer(https bool) {
	proxy := proxy.NewProxy()
	lproxy = proxy

	proxy.LoadCert(LoadCert)
	proxy.ShuntHandler(ShuntHandler)
	proxy.ReqHandler(ReqHandler)
	proxy.RespHandler(ResHandler)

	addr := fmt.Sprintf("%s:%d", config.Conf.Proxy.Host, config.Conf.Proxy.Port)
	log.Info("server listen: ", addr)

	server = &http.Server{Addr: addr, Handler: proxy}
	server.ListenAndServe()
}

func ReDeploy(b bool) error {
	if lproxy == nil {
		return nil
	}

	system.SysProxyOff()
	lproxy.CloseConnSess() //断开所有的http/https连接
	system.SysProxyOn()
	return nil
}

func ReStartServer(b bool) error {
	if server == nil {
		return nil
	}

	err := StopServer()
	go StartServer(b)

	return err
}

func StopServer() error {
	if server == nil {
		return fmt.Errorf("服务未启动")
	}

	if err := server.Close(); err != nil {
		return err
	}

	log.Info("server stopped")
	return nil
}
