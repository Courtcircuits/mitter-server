package util

import (
	"log"

	"github.com/spf13/viper"
)

func Get(key string) string {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("error while .env %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
