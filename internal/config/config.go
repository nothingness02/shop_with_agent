package cpnfig

import (
	"log"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	RunAddr string `mapstructure:"address"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"databasename"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"` // 小时为单位
}

func LoadConfig(path string) (config *Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("yaml")

	v.SetEnvPrefix("SHOP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, relying on env vars")
		} else {
			return nil, err
		}
	}
	err = v.Unmarshal(&config)
	log.Printf("DB User: %s, DB Password length: %d", config.Database.User, len(config.Database.Password))
	return
}

func (c *DatabaseConfig) BuildPostgresDSN(sslmode string) string {
	// 默认 host 和 port
	host := c.Host
	if host == "" {
		host = "localhost"
	}

	port := c.Port
	if port == "" {
		port = "5432" // PostgreSQL 默认端口
	}

	// 构建 host:port
	hostPort := host + ":" + port

	// 构建用户信息（必须提供 password，即使为空）
	userInfo := url.UserPassword(c.User, c.Password)

	u := url.URL{
		Scheme: "postgres",
		User:   userInfo,
		Host:   hostPort,
		Path:   "/" + c.DatabaseName,
	}
	q := u.Query()
	if sslmode == "" {
		sslmode = "disable"
	}
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return u.String()
}
