package env

import (
	"ethereum-api/pkg/app/config"

	"github.com/kelseyhightower/envconfig"
)

func LoadConfig() config.Config {
	config := config.Config{}
	err := envconfig.Process("ETHEREUM_API", &config)
	if err != nil {
		panic(err)
	}

	return config
}
