package main

import (
	"fmt"
	"os"
	"strings"

	"its-treason/web-test/db"
)

var requiredEnvVars = []string{
	"POSTGRES_URL",
	"ELASTICSEARCH_URL",
}

func loadConfig() (db.Config, error) {
	var missing []string
	for _, key := range requiredEnvVars {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return db.Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return db.Config{
		PostgresURL:      os.Getenv("POSTGRES_URL"),
		ElasticsearchURL: os.Getenv("ELASTICSEARCH_URL"),
	}, nil
}
