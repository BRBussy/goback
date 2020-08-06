package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	err := viper.BindEnv("ServerPort", "PORT")
	if err != nil {
		log.Fatal().Err(err).Msgf("binding environment variables to configuration keys")
	}

	// set default configuration
	viper.SetDefault("ServerPort", "8080")
	viper.SetDefault("MongoDBHosts", []string{"localhost:27017"})
	viper.SetDefault("MongoDBName", "goback")
	viper.SetDefault("PrivateKeyString", ``)

}

// Config holds configuration information for the backend
type Config struct {
	ServerPort              string
	MongoDBConnectionString string
	MongoDBHosts            []string
	MongoDBName             string
	MongoDBUsername         string
	MongoDBPassword         string
	PrivateKeyString        string
}

// GetConfig looks for and tries to parse a .toml config file with the given name
func GetConfig(configFileName *string) (*Config, error) {
	// set places to look for config file
	viper.AddConfigPath("cmd/backend/config")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	// set the name of the config file
	viper.SetConfigName(*configFileName)
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msgf("could not parse config file")
		return nil, err
	}

	// parse the config file
	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		log.Error().Err(err).Msg("unmarshalling config file")
		return nil, err
	}

	return cfg, nil
}
