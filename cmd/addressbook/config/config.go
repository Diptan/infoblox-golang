package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config contains server configuration
type Config struct {
	Db     DbConfig     `mapstructure:"db"`
	Server ServerConfig `mapstructure:"server"`
}

type DbConfig struct {
	Host string `mapstructure:"host"`
}

type ServerConfig struct {
	GatewayAddr int `mapstructure:"GatewayAddr"`
}

// ReadConfig reads config values from json config file
func Read() (Config, error) {
	vp := viper.New()
	config := Config{}

	//Read env variables and override values from config file
	vp.AutomaticEnv()
	vp.SetEnvPrefix("ADDRESSBOOK")
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vp.SetConfigName("config")
	vp.SetConfigType("json")

	//Local
	vp.AddConfigPath("config")
	// In docker
	vp.AddConfigPath("/usr/local/bin/")

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
