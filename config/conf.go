package config

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EnvironmentConfig flags
type EnvironmentConfig struct {
	CourierPort              string `env:"COURIER_PORT" mapstructure:"port"`
	CourierHost              string `env:"COURIER_HOST" mapstructure:"host"`
	CourierURL               string `env:"COURIER_WS_URL" mapstructure:"ws_url"`
	CourierDatabaseURL       string `env:"COURIER_DATABASE_URL" mapstructure:"database_url"`
	CourierDatabaseNamespace string `env:"COURIER_DATABASE_NAMESPACE" mapstructure:"database_namespace"`
	CourierJWTSecret         string `env:"COURIER_JWT_SECRET" mapstructure:"jwt_secret"`
	CourierKeySecret         string `env:"COURIER_KEY_SECRET" mapstructure:"key_secret"`
	CourierProduction        bool   `env:"COURIER_PRODUCTION" mapstructure:"production"`
}

// LoadEnvironmentConfig load config from env
func LoadEnvironmentConfig() *EnvironmentConfig {
	viper := viper.New()

	loadEnv(viper)

	envConfig := new(EnvironmentConfig)

	viper.Unmarshal(&envConfig)

	return envConfig
}

func loadDefault(viper *viper.Viper) {
	viper.SetDefault("COURIER_PRODUCTION", false)
	viper.SetDefault("COURIER_PORT", "8883")
	viper.SetDefault("COURIER_HOST", "localhost")
	viper.SetDefault("COURIER_WS_URL", "http://localhost:8883/v1/ws")
	viper.SetDefault("COURIER_DATABASE_NAMESPACE", "courier")
	viper.SetDefault("COURIER_KEY_SECRET", "kNDKzQlk8ONVfjpMKo2I7tzBxQu4nF7EwpHwrdWT6R")
}

func loadEnv(viper *viper.Viper) {
	dir, _ := os.Getwd()
	envfile := path.Join(dir, ".env")

	viper.SetEnvPrefix("courier")
	viper.SetConfigFile(envfile)

	viper.AutomaticEnv()
	viper.BindEnv("PORT")
	viper.BindEnv("HOST")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("DATABASE_NAMESPACE")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("KEY_SECRET")
	viper.BindEnv("WS_URL")
	viper.BindEnv("PRODUCTION")

	loadDefault(viper)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Info("Env file not present")
	}
}
