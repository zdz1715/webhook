package config

type Log struct {
	Level  string    `yaml:"level"`
	DIR    string    `yaml:"dir"`
	Rotate LogRotate `yaml:"rotate"`
}

type LogRotate struct {
	// MB
	MaxSize    int `yaml:"maxSize"`
	MaxBackups int `yaml:"maxBackups"`
}
