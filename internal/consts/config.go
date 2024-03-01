package consts

import (
	"net/url"
	"strings"
	"sync"
)

type Config struct {
	Env         string      `yaml:"env"`
	Mode        string      `yaml:"mode"`
	RedisConfig RedisConfig `yaml:"redis"`
	MinioConfig MinioConfig `yaml:"minio"`
	Server      Server      `yaml:"server"`
	Auth        Auth        `yaml:"auth"`
}

type Server struct {
	Address    string `yaml:"address"`
	Name       string `yaml:"name"`
	DomainPath string `yaml:"domainPath"`
	TokenKey   string `yaml:"tokenKey"`
}

func (s *Server) GetDomainScheme() string {
	parsedURL, err := url.Parse(s.DomainPath)
	if err == nil {
		return parsedURL.Scheme
	}
	return "https"
}

type RedisConfig struct {
	Default RedisConfigDetail `yaml:"default"`
}

type RedisConfigDetail struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
	Pass    string `yaml:"pass"`
}

type Auth struct {
	Login Login `yaml:"login"`
}

type Login struct {
	Expire int64 `yaml:"expire"`
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

func (config Config) IsServerDev() bool {
	return len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "server_dev") == 0
}

func (config Config) IsLocal() bool {
	return config.IsServerDev() || (len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "local") == 0)
}

func (config Config) IsStage() bool {
	return len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "daily") == 0
}

func (config Config) IsProd() bool {
	return len(config.Env) > 0 && strings.Compare(strings.ToLower(config.Env), "prod") == 0
}
