package config

import "github.com/spf13/viper"

type Conf struct {
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`

	IpLimit   int `mapstructure:"IP_LIMIT"`
	TimeForIp int `mapstructure:"IP_BLOCK_TIME_SECONDS"`

	TokenAName    string `mapstructure:"TOKEN_A_NAME"`
	TokenALimit   int    `mapstructure:"TOKEN_A_LIMIT"`
	TimeForTokenA int    `mapstructure:"TOKEN_A_BLOCK_TIME_SECONDS"`

	TokenBName    string `mapstructure:"TOKEN_B_NAME"`
	TokenBLimit   int    `mapstructure:"TOKEN_B_LIMIT"`
	TimeForTokenB int    `mapstructure:"TOKEN_B_BLOCK_TIME_SECONDS"`
}

func LoadConfig(path string) (*Conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	var cfg *Conf

	err = viper.Unmarshal(&cfg)

	if err != nil {
		panic(err)
	}

	return cfg, err
}
