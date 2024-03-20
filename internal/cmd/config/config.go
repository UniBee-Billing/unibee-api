package config

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"unibee/utility"
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
	Address     string `yaml:"address"`
	DomainPath  string `yaml:"domainPath"`
	OpenApiPath string `yaml:"openapiPath"`
	SwaggerPath string `yaml:"swaggerPath"`
	JwtKey      string `yaml:"jwtKey"`
}

func (s *Server) GetDomainScheme() string {
	parsedURL, err := url.Parse(s.DomainPath)
	if err == nil {
		return parsedURL.Scheme
	}
	return "https"
}

func (s *Server) GetServerPath() string {
	return fmt.Sprintf("%s", s.DomainPath)
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

func GetConfigInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func SetConfig(config string) {
	err := utility.UnmarshalFromJsonString(config, &instance)
	if err != nil {
		panic(err)
	}
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
