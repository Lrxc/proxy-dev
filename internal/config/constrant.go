package config

var AppName = "Proxy-Dev v1.0.1"

// 工作路径
var DataPath = "data"

// 日志路径
var LogPath = "logs/logrus.log"

const (
	PROXY_TYPE_REDIRECT = "redirect"
	PROXY_TYPE_REQMOD   = "reqmod"
	PROXY_TYPE_RESPMOD  = "respmod"
)
