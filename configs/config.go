package config

import "github.com/spf13/viper"

type conf struct {
	IpLimit      int `mapstructure:"REQUEST_LIMIT_FOR_IP"`
	TimeForIp    int `mapstructure:"TIME_FOR_IP_BLOCKED_IN_SECONDS"`
	TokenLimit   int `mapstructure:"REQUEST_LIMIT_FOR_TOKEN"`
	TimeForToken int `mapstructure:"TIME_FOR_TOKEN_BLOCKED_IN_SECONDS"`
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	var cfg *conf

	err = viper.Unmarshal(&cfg)

	if err != nil {
		panic(err)
	}

	return cfg, err
}
