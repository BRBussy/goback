package main

import (
	"flag"
	roleSyncConfig "github.com/BRBussy/goback/cmd/setup/config"
	"github.com/BRBussy/goback/pkg/logs"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/setup"
	"github.com/rs/zerolog/log"
	"time"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := roleSyncConfig.GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("error getting config")
	}

	// create new mongo db connection
	mongoDbConn, err := mongo.NewDatabaseConnection(
		config.MongoDBHosts,
		config.MongoDBUsername,
		config.MongoDBPassword,
		config.MongoDBConnectionString,
		config.MongoDBName,
		20*time.Second,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating new mongo db connection")
	}
	defer func() {
		if err := mongoDbConn.CloseConnection(); err != nil {
			log.Error().Err(err).Msg("error closing mongo db connection")
		}
	}()

	//
	// instantiate service providers
	//
	roleMongoStore := role.NewMongoStore(mongoDbConn)
	roleBasicAdmin := role.NewBasicAdmin(roleMongoStore)

	// run role sync
	log.Info().Msg("__________ Running Role Sync __________")
	setup.RoleSync(
		roleMongoStore,
		roleBasicAdmin,
	)
}
