package main

import (
	"flag"
	backendConfig "github.com/BRBussy/goback/cmd/backend/config"
	"github.com/BRBussy/goback/pkg/jsonrpc"
	"github.com/BRBussy/goback/pkg/logs"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	// get config
	config, err := backendConfig.GetConfig(configFileName)
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
	// Service Providers
	//
	userMongoStore := user.NewMongoStore(mongoDbConn)

	// create JSON-RPC HTTP server
	jsonRPCHTTPServer := jsonrpc.NewServer(
		"0.0.0.0",
		config.ServerPort,
		[]jsonrpc.RPCServerConfig{
			//
			// Public API Server
			//
			{
				Name:       "Public",
				Path:       "/api/public",
				Middleware: []mux.MiddlewareFunc{},
				ServiceProviders: []jsonrpc.ServiceProvider{
					user.NewAuthenticatedJSONRPCWrapper(userMongoStore),
				},
			},

			//
			// Authenticated API Server
			//
			{
				Name:       "Public",
				Path:       "/api/authenticated",
				Middleware: []mux.MiddlewareFunc{},
				ServiceProviders: []jsonrpc.ServiceProvider{
					user.NewAuthenticatedJSONRPCWrapper(userMongoStore),
				},
			},

			//
			// Authorised API Server
			//
			{
				Name:       "Public",
				Path:       "/api/authorised",
				Middleware: []mux.MiddlewareFunc{},
				ServiceProviders: []jsonrpc.ServiceProvider{
					user.NewAuthenticatedJSONRPCWrapper(userMongoStore),
				},
			},
		},
	)

	// start server
	go func() {
		if err := jsonRPCHTTPServer.Start(); err != nil {
			log.Error().Err(err).Msg("json rpc http api server has stopped")
		}
	}()

	// wait for interrupt signal to stop
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for s := range systemSignalsChannel {
		log.Info().Msgf("Application is shutting down.. ( %s )", s)
		return
	}
}
