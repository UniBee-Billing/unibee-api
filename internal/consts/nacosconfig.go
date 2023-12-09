package consts

import "sync"

type NacosConfig struct {
	Env         string      `yaml:"env"`
	HostPath    string      `yaml:"host_path"`
	RedisConfig RedisConfig `yaml:"redismq"`
}

type RedisConfig struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
	Pass    string `yaml:"pass"`
}

var instance *NacosConfig
var once sync.Once

// GetNacosConfigInstance 返回 Singleton 的唯一实例
func GetNacosConfigInstance() *NacosConfig {
	once.Do(func() {
		instance = &NacosConfig{}
	})
	return instance
}
