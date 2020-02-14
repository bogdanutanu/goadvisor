package config

import valid "github.com/asaskevich/govalidator"

type Config struct {
	VanguardURL string `valid:"url,required"`
}

var conf Config
