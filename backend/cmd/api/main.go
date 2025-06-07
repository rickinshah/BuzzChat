package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/RickinShah/BuzzChat/internal/data"
	"github.com/RickinShah/BuzzChat/internal/db"
	"github.com/RickinShah/BuzzChat/internal/jsonlog"
	"github.com/RickinShah/BuzzChat/internal/mailer"
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
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	port          int
	env           string
	clients       []string
	encryptionKey string
}

type application struct {
	config config
	logger *jsonlog.Logger
	mailer mailer.Mailer
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://buzzchat:buzzchat@localhost/buzzchat?sslmode=disable", "Postgres DSN")
	flag.StringVar(&cfg.redis.address, "redis-address", "redis://localhost:6379", "Redis Address")
	flag.StringVar(&cfg.encryptionKey, "encryption-key", "abcdefghijklmnopqrstuvwxyzabcdef", "Redis Address")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.gmail.com", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "no-reply@gmail.com", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "BuzzChat <no-reply@gmail.com>", "SMTP sender")

	clients := flag.String("clients", "http://localhost:5173", "Client URLs for CORS")

	flag.Parse()

	cfg.clients = strings.Split(*clients, ",")

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
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		models: data.NewModels(db.New(dbPool), redis),
	}

	app.background(func() {
		mailer.StartEmailWorker(&app.mailer, redis)
	})

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
