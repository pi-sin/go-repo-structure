package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	viper *viper.Viper
}

var configInstance *Config // singleton instance of the configuration
var singleton sync.Once    // singleton helper utility

func (c *Config) Init() {
	c.viper.SetEnvPrefix(`go-repo-struct`)
	c.viper.SetConfigType(`json`)
	c.viper.SetConfigFile(`config.json`)
	c.viper.AutomaticEnv()

	err := c.viper.ReadInConfig()

	if err != nil {
		// Handle errors reading the config file
		logrus.WithFields(logrus.Fields{
			"mod": "logger",
			"evn": "init",
		}).Error(err)

		panic(err)
	}
}

func GetConfig() *viper.Viper {
	// create an instance if not available
	singleton.Do(func() {
		configInstance = &Config{viper: viper.New()}
		configInstance.Init()
	})

	return configInstance.viper
}
