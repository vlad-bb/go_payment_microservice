package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env               string `env:"APP_ENV" envDefault:"development"`
	MongoCn           string `env:"MONGO_CN"`
	MongoDb           string `env:"MONGO_DB"`
	LoggerLevel       string `env:"LOGGER_LEVEL" envDefault:"info"`
	AuthSecret        string `env:"AUTH_SECRET"`
	TokenTTLMinutes   string `env:"TOKEN_TTL_MINUTES" envDefault:"300"`
	PyroscopeServer   string `env:"PYROSCOPE_SERVER"`
	EnableCPUProfiler string `env:"ENABLE_CPU_PROFILER" envDefault:"false"`
	InstanceName      string `env:"INSTANCE_NAME" envDefault:"dev"`
	LiqPayPublicKey   string `env:"LIQPAY_PUBLIC_KEY"`
	LiqPaySecretKey   string `env:"LIQPAY_SECRET_KEY"`
}

func NewConfigFromEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or failed to load")
	}

	cfg := &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
