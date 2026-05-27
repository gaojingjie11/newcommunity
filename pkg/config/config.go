package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Service  ServiceConfig          `mapstructure:"service"`
	MySQL    MySQLConfig            `mapstructure:"mysql"`
	Redis    RedisConfig            `mapstructure:"redis"`
	Nacos    NacosConfig            `mapstructure:"nacos"`
	RabbitMQ RabbitMQConfig         `mapstructure:"rabbitmq"`
	MinIO    MinIOConfig            `mapstructure:"minio"`
	Agent    AgentConfig            `mapstructure:"agent"`
	Gateway  GatewayConfig          `mapstructure:"gateway"`
	Raw      map[string]interface{} `mapstructure:",remain"`
}

type ServiceConfig struct {
	Name       string `mapstructure:"name"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	RegisterIP string `mapstructure:"register_ip"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=2s",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type NacosConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
}

func (c NacosConfig) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type AgentConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	LLMAPIKey  string `mapstructure:"llm_api_key"`
	LLMBaseURL string `mapstructure:"llm_base_url"`
	LLMModel   string `mapstructure:"llm_model"`
}

type GatewayConfig struct {
	Services      map[string]string `mapstructure:"services"`
	InternalToken string            `mapstructure:"internal_token"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	setDefaults(v)
	v.SetConfigFile(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("service.host", "0.0.0.0")
	v.SetDefault("mysql.host", "mysql")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.database", "smart_community")
	v.SetDefault("mysql.username", "root")
	v.SetDefault("mysql.password", "root123456")
	v.SetDefault("redis.host", "redis")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("nacos.enabled", true)
	v.SetDefault("nacos.host", "nacos")
	v.SetDefault("nacos.port", 8848)
	v.SetDefault("nacos.namespace", "public")
	v.SetDefault("nacos.group", "DEFAULT_GROUP")
	v.SetDefault("rabbitmq.host", "rabbitmq")
	v.SetDefault("rabbitmq.port", 5672)
	v.SetDefault("rabbitmq.username", "guest")
	v.SetDefault("rabbitmq.password", "guest")
	v.SetDefault("minio.endpoint", "minio:9000")
	v.SetDefault("minio.access_key", "minioadmin")
	v.SetDefault("minio.secret_key", "minioadmin")
	v.SetDefault("minio.bucket", "smart-community")
	v.SetDefault("agent.base_url", "http://agent-service:9000")
}
