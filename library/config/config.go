package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
)

type Config struct {
	Driver *viper.Viper
	watch  bool
}

func New(opts ...viper.Option) *Config {
	return &Config{
		Driver: viper.NewWithOptions(opts...),
	}
}

func (c *Config) SetDefault(key string, value interface{}) {
	c.Driver.SetDefault(key, value)
}

func (c *Config) SetWatch(b bool) {
	c.watch = b
}

// Load
// path if empty, get the default from the command line argument: --config
func (c *Config) Load(path string) (*viper.Viper, error) {
	fmt.Printf("[CONFIG] --load path: %s \n", path)
	c.Driver.SetConfigFile(path)
	// viper源码不指定文件格式的话，接下来会有多次获取文件格式，为了效率，这里指定下
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		ext = ext[1:]
		c.Driver.SetConfigType(ext)
	}
	readErr := c.Driver.ReadInConfig()
	if readErr != nil {
		return c.Driver, readErr
	}

	if c.watch {
		c.Driver.OnConfigChange(func(in fsnotify.Event) {
			readAgainErr := c.Driver.ReadInConfig()
			if readAgainErr != nil {
				fmt.Printf("[CONFIG] --watch %s msg: Read again. error: %s \n", in, readAgainErr.Error())
			} else {
				fmt.Printf("[CONFIG] --watch %s msg: Read again. \n", in)
			}
		})

		c.Driver.WatchConfig()
	}

	return c.Driver, nil
}

func (c *Config) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if c.watch {
		c.Driver.OnConfigChange(func(in fsnotify.Event) {
			if err := c.Driver.Unmarshal(rawVal, opts...); err != nil {
				fmt.Printf("[config] --watch %s msg: Unmarshal again. error: %s \n", in, err.Error())
			} else {
				fmt.Printf("[config] --watch %s msg: Unmarshal again. \n", in)
			}
		})
	}
	return c.Driver.Unmarshal(rawVal, opts...)
}
