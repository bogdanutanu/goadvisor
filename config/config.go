package config

type Config struct {
	VanguardURL string `valid:"url,required"`
}

var Cfg Config
