package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Configuration structure which hold information for configuring the datasetAPI
type Configuration struct {
	BindAddr                string        `envconfig:"BIND_ADDR"`
	DefaultMaxResults       int           `envconfig:"Default_Max_Results"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	Host                    string        `envconfig:"HOST_NAME"`
	ElasticSearchConfig     *ElasticSearchConfig
}

// ElasticSearchConfig structure which contains information to access mongo datastore
type ElasticSearchConfig struct {
	DestURL        string `envconfig:"ES_DESTINATION_URL"`
	DestIndex      string `envconfig:"ES_DESTINATION_INDEX"`
	ShowScore      bool   `envconfig:"ES_SHOW_SCORE"`
	SignedRequests bool   `envconfig:"ES_SIGNED_REQUESTS"`
}

var cfg *Configuration

// Get the application and returns the configuration structure
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		BindAddr:                ":10100",
		DefaultMaxResults:       1000,
		GracefulShutdownTimeout: 5 * time.Second,
		Host:                    "http://localhost",
		ElasticSearchConfig: &ElasticSearchConfig{
			DestURL:        "http://localhost:9200",
			DestIndex:      "courses",
			ShowScore:      false,
			SignedRequests: true,
		},
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Configuration) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
