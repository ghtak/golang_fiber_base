package core

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	LogEncoder     string `mapstructure:"LOG_ENCODER"`
	LogFileName    string `mapstructure:"LOG_FILENAME"`
	DatabaseDriver string `mapstructure:"DATABASE_DRIVER"`
	DatabaseUrl    string `mapstructure:"DATABASE_URL"`
}

func NewEnv() Env {
	env := Env{}
	var envFile *string = flag.String("cfg", ".env", "set env filename")

	flag.Parse()

	viper.SetConfigFile(*envFile)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Sprintf("Can't find file %s", *envFile), err)
	}
	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal(fmt.Sprintf("Can't Load Env %s", *envFile), err)
	}
	return env
}
