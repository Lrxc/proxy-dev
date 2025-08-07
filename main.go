package main

import (
	"embed"
	"proxy-dev/internal/config"
	"proxy-dev/internal/gui"
	"proxy-dev/internal/system"
)

//go:embed frontend/dist
var assets embed.FS

func init() {
	system.IsAlreadyRunning()
	config.InitConfig()
	config.InitLog()
}

func main() {
	gui.Gui(assets)
}
