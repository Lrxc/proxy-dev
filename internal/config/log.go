package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func InitLog() {
	file := &lumberjack.Logger{
		Filename:   LogPath,
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
		//Compress:   true, // disabled by default
	}
	//同时将日志写入文件和控制台
	writer := io.MultiWriter(file, os.Stdout)

	// 日志级别
	log.SetLevel(log.InfoLevel)
	//日志格式化
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      false,
		DisableTimestamp: true,
		TimestampFormat:  time.DateTime,
	})

	//写入文件
	log.SetOutput(writer)
}
