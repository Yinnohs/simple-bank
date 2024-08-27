package util

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	DbDriver    string `mapstructure:"DB_DRIVER"`
	DbSource    string `mapstructure:"DB_SOURCE"`
	BaseAddress string `mapstructure:"BASE_ADDRESS"`
}

func LoadConfig(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	err = viper.Unmarshal(&config)
	return
}
