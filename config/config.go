package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type configType struct {
	DB       db       `toml:"db"`
	DynamoDB dynamoDB `toml:"dynamodb,omitempty"`
	MongoDB  mongoDB  `toml:"mongo,omitempty"`
}

type db struct {
	ItemsPerPage int    `toml:"items_per_page"`
	Engine       string `toml:"engine"`
}

type dynamoDB struct {
	TableName string `toml:"table_name"`
}

type mongoDB struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

var (
	config   configType
	DB       db
	MongoDB  mongoDB
	DynamoDB dynamoDB
)

func init() {
	configFile := "./config.toml"
	_, err := toml.DecodeFile(configFile, &config)

	if err != nil {
		panic(errors.Wrapf(err, "error while reading %s", configFile))
	}

	DB = config.DB
	MongoDB = config.MongoDB
	DynamoDB = config.DynamoDB
}
