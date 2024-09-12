package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

// Config stores all configuration of the application.
// The values are read by viper from a configs file or environment variables
type Config struct {
	Server    Server    `yaml:"Server"`
	Database  Database  `yaml:"Database"`
	SnowFlake SnowFlake `yaml:"SnowFlake"`
	Redis     Redis     `yaml:"Redis"`
}

type Server struct {
	Mode   string `yaml:"Mode"`
	Host   string `yaml:"Host"`
	Port   string `yaml:"Port"`
	SLAddr string `yaml:"SLAddr"`
}

type Database struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Name     string `yaml:"Name"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
}

type SnowFlake struct {
	Node int64 `yaml:"Node"`
}

type Redis struct {
	Addr           string        `yaml:"Addr"`
	Password       string        `yaml:"Password"`
	DB             int           `yaml:"DB"`
	ExpirationTime time.Duration `yaml:"ExpirationTime"`
}

var cfg *Config

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.SetConfigName(env)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error on parsing env configuration file, %v", err)

	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("error on decoding into struct, %v", err)
	}
}

// GetConfig return the config struct
func GetConfig() *Config {
	return cfg
}
