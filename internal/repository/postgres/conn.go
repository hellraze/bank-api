package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file not found")
	}
	dns := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dns)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return pool, err
}
