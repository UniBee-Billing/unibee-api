package consts

import (
	"strings"
	"sync"
)

type Config struct {
	Env           string        `yaml:"env"`
	RedisMqConfig RedisMqConfig `yaml:"redismq"`
	MinioConfig   MinioConfig   `yaml:"minio"`
	Server        Server        `yaml:"server"`
}

type Server struct {
	Address    string `yaml:"address"`
	Name       string `yaml:"name"`
	DomainPath string `yaml:"domainPath"`
}

type RedisMqConfig struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
	Pass    string `yaml:"pass"`
}

type MinioConfig struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	BucketName string `yaml:"bucketName"`
	Domain     string `yaml:"domain"`
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

func (config Config) IsLocal() bool {
	return len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "local") == 0
}

func (config Config) IsStage() bool {
	return len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "daily") == 0
}

func (config Config) IsProd() bool {
	return len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "prod") == 0
}
