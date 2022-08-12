package platform

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	config *viper.Viper

	loadConfigOnce sync.Once
)

func GetConfig() *viper.Viper {
	return config
}

func LoadConfig() (rErr error) {
	loadConfigOnce.Do(func() {
		config = viper.New()

		config.AddConfigPath(".")

		config.SetConfigName(".dddoc")

		config.SetConfigType("yaml")

		rErr = config.ReadInConfig()
	})

	return
}
