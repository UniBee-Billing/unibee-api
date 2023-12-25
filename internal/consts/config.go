package consts

import "sync"

type Config struct {
	Env         string      `yaml:"env"`
	HostPath    string      `yaml:"host_path"`
	RedisConfig RedisConfig `yaml:"redismq"`
	Server      Server      `yaml:"server"`
}

type Server struct {
	Address string `yaml:"address"`
}

type RedisConfig struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
	Pass    string `yaml:"pass"`
}

var instance *Config
var once sync.Once

// GetConfigInstance 返回 Singleton 的唯一实例
func GetConfigInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}
