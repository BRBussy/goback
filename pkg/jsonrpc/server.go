package jsonrpc

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	gorillaJSON "github.com/gorilla/rpc/v2/json2"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Server struct {
	host   string
	port   string
	router *mux.Router
}

type RPCServerConfig struct {
	Name             string
	Path             string
	Middleware       []mux.MiddlewareFunc
	ServiceProviders []ServiceProvider
}

func NewServer(
	host string,
	port string,
	rpcServerConfigs []RPCServerConfig,
) *Server {
	// create a new Server
	newServer := new(Server)
	newServer.host = host
	newServer.port = port

	newServer.router = mux.NewRouter()
	newServer.router.Use(preFlightAndCORSHandler)

	for _, rpcServerConfig := range rpcServerConfigs {
		log.Info().Msg(fmt.Sprintf(
			"Start %s RPC API Server on path %s",
			rpcServerConfig.Name,
			rpcServerConfig.Path,
		))

		// create new gorilla rpc server and register JSON codec
		rpcServer := rpc.NewServer()
		rpcServer.RegisterCodec(gorillaJSON.NewCodec(), "application/json")

		// register each service provider with the rpc server
		for _, serviceProvider := range rpcServerConfig.ServiceProviders {
			log.Info().Msg("	Registering: " + serviceProvider.ServiceProviderName())
			if err := rpcServer.RegisterService(serviceProvider, serviceProvider.ServiceProviderName()); err != nil {
				log.Fatal().Err(err).Msg("could not register: " + serviceProvider.ServiceProviderName())
			}
		}

		// create api router
		apiRouter := mux.NewRouter()

		// register any supplied middleware
		if rpcServerConfig.Middleware != nil {
			apiRouter.Use(rpcServerConfig.Middleware...)
		}

		// put a handler function
		apiRouter.HandleFunc("/", rpcServer.ServeHTTP)

		// mount this individual RPC Server and the given path
		newServer.router.Handle(rpcServerConfig.Path, apiRouter)
	}

	return newServer
}

func (s *Server) Start() error {
	log.Info().Msg("starting http json rpc api server on: " + s.host + ":" + s.port)
	return http.ListenAndServe(s.host+":"+s.port, s.router)
}

func preFlightAndCORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		}
	})
}
