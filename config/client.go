package config

import "time"

type Client struct {
	Timeout time.Duration `yaml:"timeout"`
	//RetryCount     int           `yaml:"retryCount"`
	//RetrySleepTime time.Duration `yaml:"retrySleepTime"`
}
