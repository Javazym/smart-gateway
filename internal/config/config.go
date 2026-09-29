package config

import (
	"github.com/spf13/viper"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 网关整体配置。
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Jwt      JwtConfig      `yaml:"jwt"`
	Security SecurityConfig `yaml:"security"`
}

// SecurityConfig 安全相关配置。
type SecurityConfig struct {
	// SSRFProtection 是否拦截内网/回环/链路本地等上游地址。
	// 本网关定位为内网/单机反向代理，默认 false（放行内网上游）；
	// 对外暴露且上游不可信时置 true 以阻断 SSRF。
	SSRFProtection bool `yaml:"ssrf_protection"`
}

// ServerConfig HTTP 服务相关配置。
type ServerConfig struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// JwtConfig JWT 鉴权相关配置。
type JwtConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

// Route 描述一条转发规则：把匹配的请求转发到指定上游地址。

// Mysql配置
type MysqlConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime string `mapstructure:"connMaxLifetime"`
}

// redis配置
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`     // 127.0.0.1:6379
	DB           int    `mapstructure:"db"`       // 使用哪个db，默认0
	PoolSize     int    `mapstructure:"poolSize"` // 连接池大小
	MinIdleConns int    `mapstructure:"minIdleConns"`
}

// Load 从 YAML 文件读取并解析配置。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// 读取mysql配置
func LoadMysqlConfig() (*MysqlConfig, error) {
	viper.SetConfigFile("configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg MysqlConfig
	err = viper.UnmarshalKey("mysql", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// 读取redis配置
func LoadRedisConfig() (*RedisConfig, error) {
	viper.SetConfigFile("configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg RedisConfig
	err = viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
