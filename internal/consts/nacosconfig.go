package consts

import "sync"

type NacosConfig struct {
	Env      string `yaml:"env"`
	HostPath string `yaml:"host_path"`
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
