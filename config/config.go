package config

import (
	"github.com/spf13/pflag"
	"github.com/zdz1715/webhook/pkg/config"
	"time"
)

type Config struct {
	Application Application        `yaml:"application"`
	Webhooks    map[string]Webhook `yaml:"webhooks"`
	Client      Client             `yaml:"client"`
	Log         Log                `yaml:"log"`
}

var (
	cfg  *Config
	path string
)

// TODO: 增加超时重试功能
func init() {
	// main.go pflag.Parse()
	pflag.StringVarP(&path, "config", "c", "config.yaml", "choose config file.")
}

func Init() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = defaultConfig(cfg)
	c := config.New()
	c.SetWatch(true)
	_, err := c.Load(path)
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func defaultConfig(cfg *Config) *Config {
	cfg = new(Config)
	cfg.Application = Application{
		Port:         8000,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	cfg.Client = Client{
		Timeout: 30 * time.Second,
		//RetryCount: 3,
	}

	cfg.Log = Log{
		DIR: "./logs",
		Rotate: LogRotate{
			MaxSize:    0,
			MaxBackups: 0,
		},
	}

	return cfg
}
