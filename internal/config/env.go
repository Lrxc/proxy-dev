package config

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"proxy-dev/internal/util"
)

type Config struct {
	System System
	Proxy  Proxy
	Rule   []Rule
}

type System struct {
	MinExit   bool
	Https     bool
	AutoProxy bool `mapstructure:"auto_proxy" yaml:"auto_proxy"`
}

type Proxy struct {
	Host string
	Port int
}

type Rule struct {
	Enable bool
	Type   string

	Surl string
	Turl string

	Sdata string
	Tdata string
}

const confname = "conf.yml"

// 初始化环境参数
func InitConfig() {
	exit := util.FileExist(confname)
	if !exit {
		json := &Config{
			Proxy: Proxy{Host: "127.0.0.1", Port: 10086},
			Rule: []Rule{
				{Enable: true, Type: PROXY_TYPE_REDIRECT, Surl: "https://www.baidu.com", Turl: "https://www.bing.com"},
				{Enable: true, Type: PROXY_TYPE_RESPMOD, Surl: "https://api.thecatapi.com/v1/images/search?limit=", Sdata: "url", Tdata: "***"},
			},
		}
		//写入默认配置文件
		WriteConf(json)
	}

	viper.AddConfigPath("../../") //配置文件路径
	viper.SetConfigFile("conf.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	ReadConf("init")

	// 监听配置文件的变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		ReadConf("reload")
	})
	viper.WatchConfig()
}

func ReadConf(msg string) {
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	log.Info("config ", msg)
}

func WriteJson(msg string) error {
	err := json.Unmarshal([]byte(msg), &Conf.Rule)
	if err != nil {
		return err
	}
	return WriteConf(Conf)
}

func WriteConf(conf *Config) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	join, _ := filepath.Abs(confname)
	return os.WriteFile(join, data, 0644)
}
