package api

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/methods/go-methods-lib/log"
	"github.com/methods/go-methods-lib/server"
	"github.com/ofs/alpha-search-api/config"
)

var (
	httpServer   *server.Server
	serverErrors chan error
)

// API provides an interface for the routes
type API interface {
	CreateSearchAPI(string, *mux.Router) *SearchAPI
}

// SearchAPI manages stored data
type SearchAPI struct {
	DefaultMaxResults int
	Elasticsearch     Elasticsearcher
	Host              string
	Index             string
	Router            *mux.Router
	ShowScore         bool
}

// CreateSearchAPI manages all the routes configured to API
func CreateSearchAPI(cfg config.Configuration, elasticsearch Elasticsearcher, errorChan chan error) {
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
func Routes(cfg config.Configuration, elasticsearch Elasticsearcher, router *mux.Router) *SearchAPI {

	host := cfg.Host + cfg.BindAddr

	api := SearchAPI{
		DefaultMaxResults: cfg.DefaultMaxResults,
		Elasticsearch:     elasticsearch,
		Host:              host,
		Index:             cfg.ElasticSearchConfig.DestIndex,
		Router:            router,
		ShowScore:         cfg.ElasticSearchConfig.ShowScore,
	}

	api.Router.HandleFunc("/search/courses", api.SearchCourses).Methods("GET")
	api.Router.HandleFunc("/search/institution-courses", api.SearchInstitutionCourses).Methods("GET")
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
