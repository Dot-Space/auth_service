package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	DB             DBConfig
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD" env-default:"postgres"`
	DBName   string `env:"DBNAME" env-default:"postgres"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("Config path is empty!")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file does not exists!")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Config file is empty!")
	}

	if err := cleanenv.ReadEnv(&cfg.DB); err != nil {
		panic("Environment variables are empty!")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "Path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
