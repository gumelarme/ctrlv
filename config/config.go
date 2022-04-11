package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	ItemsPerPage int `toml:"items_per_page"`
	DB           DB  `toml:"db"`
}

type DB struct {
	TableName string `toml:"table_name"`
}

var Conf Config

func init() {
	configFile := "./config.toml"
	_, err := toml.DecodeFile(configFile, &Conf)

	if err != nil {
		panic(errors.Wrapf(err, "error while reading %s", configFile))
	}
}
