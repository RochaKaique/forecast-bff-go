package server

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func NewConfiguration() (*viper.Viper, error) {
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../")
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event){
		fmt.Println("Configuration changed:", e.Name)
	})

	return viper.GetViper(), nil
}