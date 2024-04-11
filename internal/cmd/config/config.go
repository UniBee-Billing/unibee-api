package config

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"unibee/utility"
)

type Config struct {
	Env         string      `json:"env" yaml:"env"`
	Mode        string      `json:"mode" yaml:"mode"`
	RedisConfig RedisConfig `json:"redis" yaml:"redis"`
	MinioConfig MinioConfig `json:"minio" yaml:"minio"`
	Server      Server      `json:"server" yaml:"server"`
	Auth        Auth        `json:"auth" yaml:"auth"`
	VatConfig   VatConfig   `json:"vatConfig" yaml:"vatConfig"`
}

type Server struct {
	Address     string `json:"address" yaml:"address"`
	DomainPath  string `json:"domainPath" yaml:"domainPath"`
	OpenApiPath string `json:"openapiPath" yaml:"openapiPath"`
	SwaggerPath string `json:"swaggerPath" yaml:"swaggerPath"`
	JwtKey      string `json:"jwtKey" yaml:"jwtKey"`
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
	Default RedisConfigDetail `json:"default" yaml:"default"`
}

type RedisConfigDetail struct {
	Address string `json:"address" yaml:"address"`
	DB      int    `json:"db" yaml:"db"`
	Pass    string `json:"pass" yaml:"pass"`
}

type Auth struct {
	Login Login `json:"login" yaml:"login"`
}

type Login struct {
	Expire int64 `json:"expire" yaml:"expire"`
}

type MinioConfig struct {
	Endpoint   string `json:"endpoint" yaml:"endpoint"`
	AccessKey  string `json:"accessKey" yaml:"accessKey"`
	SecretKey  string `json:"secretKey" yaml:"secretKey"`
	BucketName string `json:"bucketName" yaml:"bucketName"`
	Domain     string `json:"domain" yaml:"domain"`
}

type VatConfig struct {
	NonEuEnable                   string `json:"nonEuEnable" yaml:"nonEuEnable"`
	NumberUnExemptionCountryCodes string `json:"numberUnExemptionCountryCodes" yaml:"numberUnExemptionCountryCodes"`
}

var instance *Config
var once sync.Once

func GetConfigInstance() *Config {
	once.Do(func() {
		if instance == nil {
			instance = &Config{}
		}
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
