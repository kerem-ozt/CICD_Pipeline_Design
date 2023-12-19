package services

import (
	"github.com/kerem-ozt/GoodBlast_API/models"
	"github.com/spf13/viper"
)

var Config *models.EnvConfig

func LoadConfig() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetDefault("SERVER_PORT", "3002")
	v.SetDefault("MODE", "debug")
	v.SetConfigType("dotenv")
	v.SetConfigName(".env")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}

	if err := Config.Validate(); err != nil {
		panic(err)
	}
}
