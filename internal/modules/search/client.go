package search

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog"
)

func NewClient(url string, log zerolog.Logger) (*elasticsearch.Client, error) {
	if url == "" {
		log.Info().Msg("elasticsearch not configured — search features will use database fallback")
		return nil, nil
	}

	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	// Check connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("error connecting to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	log.Info().Str("url", url).Msg("connected to elasticsearch")
	return client, nil
}
