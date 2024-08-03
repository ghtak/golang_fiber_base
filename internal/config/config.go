package config

import (
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"os"
	"strings"
)

type Config struct {
	AppPort     string
	LogLevel    string
	LogEncoder  string // json, console
	LogFileName string
}

func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func NewConfig() Config {
	cfg := flag.String("cfg", "env", "env|[profile].env")
	flag.Parse()
	if strings.HasSuffix(*cfg, ".env") {
		_ = godotenv.Load(*cfg)
	}
	return Config{
		AppPort:    Getenv("APP_PORT", "3004"),
		LogLevel:   Getenv("LOG_LEVEL", ""),
		LogEncoder: Getenv("LOG_ENCODER", "console"),
		//LogFileName: "./logs/app.log",
		LogFileName: Getenv("LOG_FILE_NAME", ""),
	}
}

func Module() fx.Option {
	return fx.Module(
		"ModuleConfig",
		fx.Provide(NewConfig),
	)
}
