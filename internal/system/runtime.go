package system

import (
	"fyne.io/fyne/v2"
	"net"
	"os"
	"path/filepath"
	"proxy-dev/internal/config"
	"proxy-dev/internal/util"
)

func IsAlreadyRunning() {
	lockFile := filepath.Join(os.TempDir(), "proxy-dev.lock")

	exist := util.FileExist(lockFile)
	if exist {
		err := os.Remove(lockFile)
		if err != nil {
			sendActivateSignal()
			os.Exit(1)
		}
	}

	file, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		os.Exit(1)
	}

	go setupIPCServer()

	//保证file不被释放(程序单实例)
	go func() {
		c := make(chan []byte)
		b := <-c
		file.Write(b)
	}()
}

func sendActivateSignal() {
	conn, err := net.Dial("tcp", "127.0.0.1:51951")
	if err != nil {
		return
	}
	defer conn.Close()
	conn.Write([]byte("activate")) // 发送激活信号
}

// 监听激活信号
func setupIPCServer() {
	l, err := net.Listen("tcp", "127.0.0.1:51951")
	if err != nil {
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()

		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		if string(buf[:n]) == "activate" {
			fyne.Do(func() {
				config.AppWindow.Show()
			})
		}
	}
}
