package config

import "time"

type Application struct {
	Host         string        `yaml:"host"`
	Port         int64         `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}
