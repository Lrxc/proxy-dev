package main

import (
	"proxy-dev/internal/config"
	"proxy-dev/internal/gui"
	"proxy-dev/internal/system"
)

func init() {
	system.IsAlreadyRunning()
	config.InitLog()
	config.InitConfig()
}

func main() {
	gui.Gui()
}
