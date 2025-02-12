package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	IsProdEnv bool
	App       AppConfig      `mapstructure:"app"`
	Database  DatabaseConfig `mapstructure:"database"`
	Tracing   OtelTracing    `mapstructure:"tracing"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	Name           string        `mapstructure:"name"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	MaxConnections int           `mapstructure:"max_connections"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

type OtelTracing struct {
	Endpoint string `mapstructure:"endpoint"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./"
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct: %w", err)
	}

	config.IsProdEnv = config.App.Env == "production"

	return config, nil
}

func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("load config failed: ", err)
	}

	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
	)
}
