package db

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the database connections for the application.
type DB struct {
	Postgres      *pgxpool.Pool
	Elasticsearch *elasticsearch.Client
}

// Config holds the configuration for database connections.
type Config struct {
	PostgresURL      string
	ElasticsearchURL string
}

// New creates a new DB instance with connections to PostgreSQL and Elasticsearch.
func New(ctx context.Context, cfg Config) (*DB, error) {
	pg, err := connectPostgres(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	es, err := connectElasticsearch(cfg.ElasticsearchURL)
	if err != nil {
		pg.Close()
		return nil, fmt.Errorf("elasticsearch connection failed: %w", err)
	}

	return &DB{
		Postgres:      pg,
		Elasticsearch: es,
	}, nil
}

// Close closes all database connections.
func (d *DB) Close() {
	if d.Postgres != nil {
		d.Postgres.Close()
		log.Info("PostgreSQL connection closed")
	}
}

func connectPostgres(ctx context.Context, connURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Info("Connected to PostgreSQL")
	return pool, nil
}

func connectElasticsearch(address string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create elasticsearch client: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error response: %s", res.String())
	}

	log.Info("Connected to Elasticsearch")
	return client, nil
}
