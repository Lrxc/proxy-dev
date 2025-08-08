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
	Type   string `json:"type,omitempty"`

	Surl string `json:"surl,omitempty"`
	Turl string `json:"turl,omitempty"`

	Sdata string `json:"sdata,omitempty"`
	Tdata string `json:"tdata,omitempty"`
}

const confname = "conf.yml"

// 初始化环境参数
func InitConfig() {
	exit := util.FileExist(confname)
	if !exit {
		json := &Config{
			System: System{MinExit: true},
			Proxy:  Proxy{Host: "127.0.0.1", Port: 10086},
			Rule: []Rule{
				{Enable: true, Type: PROXY_TYPE_REDIRECT, Surl: "https://www.baidu.com", Turl: "https://www.bing.com"},
				{Enable: true, Type: PROXY_TYPE_RESPMOD, Surl: "https://www.test.com", Sdata: "name", Tdata: "alias"},
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
	var rule []Rule
	err := json.Unmarshal([]byte(msg), &rule)
	if err != nil {
		return err
	}

	Conf.Rule = rule
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
