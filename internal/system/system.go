package system

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wzshiming/sysproxy"
	//_ "github.com/wzshiming/sysproxy"
	"os"
	"os/signal"
	"proxy-dev/internal/config"
	"sync"
	"syscall"
	_ "unsafe"
)

var (
	Once    sync.Once
	SigChan = make(chan os.Signal, 10)
)

////go:linkname set github.com/wzshiming/sysproxy.set
//func set(key string, typ string, value string) error

func SysProxyOn() error {
	// 启动时设置系统代理
	addr := fmt.Sprintf("%s:%d", config.Conf.Proxy.Host, config.Conf.Proxy.Port)
	err := sysproxy.OnHTTP(addr)
	if config.Conf.System.Https {
		err = sysproxy.OnHTTPS(addr)
	}
	if err != nil {
		log.Errorf("system proxy err: %v", err)
		return err
	}
	log.Warn("system proxy on: ", addr)

	go ExitFunc()
	return nil
}

func SysProxyOff() error {
	err := sysproxy.OffHTTP()
	err = sysproxy.OffHTTPS()
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
