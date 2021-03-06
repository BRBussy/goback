package main

import (
	"crypto/rsa"
	"flag"
	backendConfig "github.com/BRBussy/goback/cmd/backend/config"
	"github.com/BRBussy/goback/pkg/jsonrpc"
	"github.com/BRBussy/goback/pkg/key"
	"github.com/BRBussy/goback/pkg/logs"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/security/authentication"
	"github.com/BRBussy/goback/pkg/security/authorisation"
	"github.com/BRBussy/goback/pkg/security/jwt"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

var configFileName = flag.String("config-file-name", "config", "specify config file")

// TODO: CORS check! --> get the server to work with the preflight thing

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

	// parse or generate RSA key pair
	var rsaPrivateKey *rsa.PrivateKey
	if config.PrivateKeyString == "" {
		pk, err := key.GenerateRSAPrivateKey()
		if err != nil {
			log.Fatal().Err(err).Msg("error generating rsa private key")
		}
		rsaPrivateKey = pk
	} else {
		pk, err := key.ParseRSAPrivateKeyFromString(config.PrivateKeyString)
		if err != nil {
			log.Fatal().Err(err).Msg("error parsing rsa private key")
		}
		rsaPrivateKey = pk
	}

	//
	// instantiate service providers
	//
	roleMongoStore := role.NewMongoStore(mongoDbConn)

	userMongoStore := user.NewMongoStore(mongoDbConn)
	userBasicAdmin := user.NewBasicAdmin(
		userMongoStore,
		roleMongoStore,
	)

	jwtBasicGenerator := jwt.NewBasicGenerator(rsaPrivateKey)
	jwtBasicValidator := jwt.NewBasicValidator(rsaPrivateKey)

	basicAuthoriser := authorisation.NewBasicAuthoriser(
		userMongoStore,
		roleMongoStore,
	)

	basicAuthenticator := authentication.NewBasicAuthenticator(
		userMongoStore,
		jwtBasicGenerator,
		jwtBasicValidator,
	)

	//
	// instantiate middleware
	//
	authenticationMiddleware := authentication.NewMiddleware(
		basicAuthenticator,
	)
	authorisationMiddleware := authorisation.NewMiddleware(
		basicAuthoriser,
		userMongoStore,
		roleMongoStore,
	)

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
					authentication.NewAuthenticatorAuthorisedJSONRPCWrapper(basicAuthenticator),
				},
			},

			//
			// Authenticated API Server
			//
			{
				Name: "Public",
				Path: "/api/authenticated",
				Middleware: []mux.MiddlewareFunc{
					authenticationMiddleware.Apply,
				},
				ServiceProviders: []jsonrpc.ServiceProvider{},
			},

			//
			// Authorised API Server
			//
			{
				Name: "Public",
				Path: "/api/authorised",
				Middleware: []mux.MiddlewareFunc{
					authenticationMiddleware.Apply,
					authorisationMiddleware.Apply,
				},
				ServiceProviders: []jsonrpc.ServiceProvider{
					user.NewStoreAuthorisedJSONRPCWrapper(userMongoStore),
					user.NewAdminAuthorisedJSONRPCWrapper(userBasicAdmin),
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
