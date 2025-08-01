package system

import (
	"fmt"
	"github.com/Trisia/gosysproxy"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"proxy-dev/internal/config"
	"sync"
	"syscall"
)

var (
	Once    sync.Once
	SigChan = make(chan os.Signal, 10)
)

func SysProxyOn() error {
	// 启动时设置系统代理
	addr := fmt.Sprintf("%s:%d", config.Conf.Proxy.Host, config.Conf.Proxy.Port)
	err := gosysproxy.SetGlobalProxy(addr)
	if err != nil {
		log.Errorf("system proxy err: %v", err)
		return err
	}
	log.Warn("system proxy on: ", addr)

	go ExitFunc()
	return nil
}

func SysProxyOff() error {
	err := gosysproxy.Off()
	log.Warn("system proxy off: ", err == nil)
	return err
}

func ExitFunc() {
	Once.Do(func() {
		signal.Notify(SigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
		s := <-SigChan
		log.Warn("exit: ", s)

		SysProxyOff()
		os.Exit(0)
	})
}
