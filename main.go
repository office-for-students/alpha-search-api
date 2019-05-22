package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/methods/go-methods-lib/log"
	"github.com/ofs/alpha-search-api/api"
	"github.com/ofs/alpha-search-api/config"
	"github.com/ofs/alpha-search-api/elasticsearch"
)

func main() {
	log.Namespace = "alpha-search-api"

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		log.ErrorC("errored getting configuration", err, log.Data{"config": cfg})
		os.Exit(1)
	}

	log.Info("configuration on startup", log.Data{"config": cfg})

	elasticClient := http.Client{}
	elasticsearch := elasticsearch.NewElasticSearchAPI(elasticClient, cfg.ElasticSearchConfig.DestURL, cfg.ElasticSearchConfig.SignedRequests)

	// Check elastic search connection can be made
	_, status, err := elasticsearch.CallElastic(context.Background(), cfg.ElasticSearchConfig.DestURL, "GET", nil)
	if err != nil {
		log.ErrorC("failed to start up, unable to connect to elastic search instance", err, log.Data{"http_status": status})
		os.Exit(1)
	}

	apiErrors := make(chan error, 1)

	api.CreateSearchAPI(*cfg, elasticsearch, apiErrors)

	// Gracefully shutdown the application closing any open resources.
	gracefulShutdown := func() {
		start := time.Now()
		log.Info("shutdown with timeout in seconds", log.Data{"timeout": cfg.GracefulShutdownTimeout})
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)

		// stop any incoming requests before closing any outbound connections
		api.Close(ctx)

		// TODO close connection to database

		log.Info("shutdown complete", log.Data{"shutdown_duration": time.Since(start)})

		cancel()
		os.Exit(1)
	}

	for {
		select {
		case err := <-apiErrors:
			log.Info("api error received", log.Data{"api_error": err})
			gracefulShutdown()
		case signal := <-signals:
			log.Info("os signal received", log.Data{"os_signal": signal})
			gracefulShutdown()
		}
	}
}
