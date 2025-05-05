package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/RickinShah/BuzzChat/internal/data"
	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/RickinShah/BuzzChat/internal/jsonlog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type config struct {
	db struct {
		dsn string
	}
	redis struct {
		address string
	}
	port int
	env  string
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://buzzchat:buzzchat@localhost/buzzchat?sslmode=disable", "Postgres DSN")
	flag.StringVar(&cfg.redis.address, "redis-address", "redis://localhost:6379", "Redis Address")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	dbPool, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer dbPool.Close()

	logger.PrintInfo("database connection pool established", nil)

	redis, err := openRedis(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer redis.Close()

	logger.PrintInfo("valkey redis connection established", nil)

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db.New(dbPool), redis),
	}

	if err = app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func openRedis(cfg config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.redis.address)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return client, nil
}
