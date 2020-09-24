package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetConfig location
func SetConfig(p string) {
	viper.SetConfigName("App")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(p)
	viper.AddConfigPath("./configurations")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config error: ", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Warn("Config file changed:", e.Name)
	})
}
