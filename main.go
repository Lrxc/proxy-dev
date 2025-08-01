package main

import (
	"proxy-dev/internal/config"
	"proxy-dev/internal/gui"
	"proxy-dev/internal/system"
)

func init() {
	system.IsAlreadyRunning()
	config.InitConfig()
	config.InitLog()
}

func main() {
	gui.Gui()
}
