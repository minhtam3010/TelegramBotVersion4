package singleton

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost        string `mapstructure:"SERVER_HOST"`
	Port              string `mapstructure:"SERVER_PORT"`
	ReadTimeout       int    `mapstructure:"SERVER_READ_TIMEOUT"`
	ReadHeaderTimeout int    `mapstructure:"SERVER_READ_HEADER_TIMEOUT"`
	WriteTimeout      int    `mapstructure:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout       int    `mapstructure:"SERVER_IDLE_TIMEOUT"`
	MaxHeaderBytes    int    `mapstructure:"SERVER_MAX_HEADER_BYTES"`

	Type         string `mapstructure:"DATABASE_TYPE"`
	User         string `mapstructure:"DATABASE_USER"`
	Password     string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseHost string `mapstructure:"DATABASE_HOST"`
	Name         string `mapstructure:"DATABASE_NAME"`
	DatabasePort string `mapstructure:"DATABASE_PORT"`
	SSLMode      string `mapstructure:"DATABASE_SSL_MODE"`
	CACERTBASE64 string `mapstructure:"DATABASE_CACERTBASE64"`

	CookieDomain   string `mapstructure:"COOKIE_DOMAIN"`
	CookieHttpOnly bool   `mapstructure:"COOKIE_HTTP_ONLY"`
	CookieSecure   bool   `mapstructure:"COOKIE_SECURE"`

	HTTPSameSite int    `mapstructure:"HTTP_SAME_SITE"`
	HTTPDomain   string `mapstructure:"HTTP_DOMAIN"`
	JWTKey       string `mapstructure:"JWT_KEY"`
}

var Cfg *Config

func InitConfig(path string) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}
