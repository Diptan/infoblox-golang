package config

import "github.com/spf13/viper"

// Config contains server configuration
type Config struct {
	Db     DbConfig     `mapstructure:"address"`
	Server ServerConfig `mapstructure:"server"`
}

type DbConfig struct {
	Address string `mapstructure:"address"`
}

type ServerConfig struct {
	HttpServerPort        int `mapstructure:"HttpServerPort"`
	GrpcServerPort        int `mapstructure:"GrpcServerPort"`
	GrpcGatewayServerPort int `mapstructure:"GrpcGatewayServerPort"`
}

// ReadConfig reads config values from json config file
func Read() (Config, error) {
	vp := viper.New()
	config := Config{}

	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("./config")

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)

	return config, nil
}
