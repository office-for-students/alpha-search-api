package api

import (
	"context"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"github.com/gorilla/mux"
	"github.com/ofs/alpha-search-api/config"
	"github.com/ofs/alpha-search-api/elasticsearch"
	"github.com/ofs/alpha-search-api/handlers"
)

var (
	httpServer   *server.Server
	serverErrors chan error
)

// SearchAPI manages stored data
type SearchAPI struct {
	Router *mux.Router
}

// CreateSearchAPI manages all the routes configured to API
func CreateSearchAPI(cfg config.Configuration, elasticsearch elasticsearch.Elasticsearcher, errorChan chan error) {
	router := mux.NewRouter()
	Routes(cfg, elasticsearch, router)

	httpServer = server.New(cfg.BindAddr, router)

	// Disable this here to allow main to manage graceful shutdown of the entire app.
	httpServer.HandleOSSignals = false

	go func() {
		log.Info("Starting search API...", nil)
		if err := httpServer.ListenAndServe(); err != nil {
			log.ErrorC("search API http server returned error", err, nil)
			errorChan <- err
		}
	}()
}

// Routes represents a list of endpoints that exist with this api
func Routes(cfg config.Configuration, elasticsearch elasticsearch.Elasticsearcher, router *mux.Router) *SearchAPI {

	api := SearchAPI{
		Router: router,
	}

	host := cfg.Host + cfg.BindAddr

	routeConfig := handlers.API{Elasticsearch: elasticsearch, Index: cfg.ElasticSearchConfig.DestIndex, Host: host, DefaultMaxResults: cfg.DefaultMaxResults}
	api.Router.HandleFunc("/search", routeConfig.AllSearch).Methods("GET")
	return &api
}

// Close represents the graceful shutting down of the http server
func Close(ctx context.Context) error {
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Info("graceful shutdown of http server complete", nil)
	return nil
}
