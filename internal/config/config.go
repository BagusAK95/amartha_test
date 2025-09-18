package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Application    ApplicationConfig
	Postgres       PostgresConfig
	Mail           MailConfig
	ContextTimeout time.Duration
}

type ApplicationConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Env  string `mapstructure:"APP_ENV"`
	Port int    `mapstructure:"APP_PORT"`
	Url  string `mapstructure:"APP_URL"`
}

type PostgresConfig struct {
	MasterHost         string        `mapstructure:"POSTGRES_MASTER_HOST"`
	MasterUsername     string        `mapstructure:"POSTGRES_MASTER_USERNAME"`
	MasterPassword     string        `mapstructure:"POSTGRES_MASTER_PASSWORD"`
	MasterPort         int           `mapstructure:"POSTGRES_MASTER_PORT"`
	MasterSSLMode      string        `mapstructure:"POSTGRES_MASTER_SSL_MODE"`
	SlaveHost          string        `mapstructure:"POSTGRES_SLAVE_HOST"`
	SlaveUsername      string        `mapstructure:"POSTGRES_SLAVE_USERNAME"`
	SlavePassword      string        `mapstructure:"POSTGRES_SLAVE_PASSWORD"`
	SlavePort          int           `mapstructure:"POSTGRES_SLAVE_PORT"`
	SlaveSSLMode       string        `mapstructure:"POSTGRES_SLAVE_SSL_MODE"`
	Database           string        `mapstructure:"POSTGRES_DATABASE"`
	Timezone           string        `mapstructure:"POSTGRES_TIMEZONE"`
	MaxOpenConnections int           `mapstructure:"POSTGRES_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `mapstructure:"POSTGRES_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `mapstructure:"POSTGRES_CONN_MAX_LIFETIME"`
}

// Global config
var APP_URL string

type MailConfig struct {
	Host     string `mapstructure:"MAIL_HOST"`
	Port     int    `mapstructure:"MAIL_PORT"`
	Username string `mapstructure:"MAIL_USERNAME"`
	Password string `mapstructure:"MAIL_PASSWORD"`
}

func Load() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	setDefaultConfig()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// Unmarshal each section explicitly
	if err = viper.Unmarshal(&config.Application); err != nil {
		return
	}
	if err = viper.Unmarshal(&config.Postgres); err != nil {
		return
	}
	if err = viper.Unmarshal(&config.Mail); err != nil {
		return
	}

	config.ContextTimeout, err = time.ParseDuration(viper.GetString("CONTEXT_TIMEOUT") + "s")

	APP_URL = viper.GetString("APP_URL")

	return
}

func setDefaultConfig() {
	viper.SetDefault("CONTEXT_TIMEOUT", 5)

	viper.SetDefault("POSTGRES_TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("POSTGRES_MAX_OPEN_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_CONN_MAX_LIFETIME", 300)
}
