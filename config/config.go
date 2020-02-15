package config

type Config struct {
	VanguardURL string `mapstructure:"vanguard_url" valid:"required"`
}

var Cfg Config
