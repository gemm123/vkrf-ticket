package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func InitConnPool() *pgxpool.Pool {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", DBUser(), DBPassword(), DBHost(), DBPort(), DBName())

	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = 30 * 60 * 1000
	dbConfig.MaxConnIdleTime = 10 * 60 * 1000

	poolConfig, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	err = poolConfig.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	return poolConfig

}
