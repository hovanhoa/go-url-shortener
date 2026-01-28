package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hovanhoa/go-url-shortener/config"
	"github.com/hovanhoa/go-url-shortener/internal/handler"
	"github.com/hovanhoa/go-url-shortener/internal/service"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
	"github.com/hovanhoa/go-url-shortener/pkg/logger"
	"github.com/hovanhoa/go-url-shortener/pkg/otel"
	"github.com/hovanhoa/go-url-shortener/pkg/snowflake"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	_ "github.com/lib/pq" // Import the pq driver
)

func Init() {
	cfg := config.GetConfig()

	// Initialize structured JSON logger
	logger.Init("go-url-shortener")
	slog.Info("Starting application", "port", cfg.Server.Port)

	// Initialize OpenTelemetry
	shutdown, err := otel.Init("go-url-shortener", "")
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
		// Continue without OTel if initialization fails
	} else {
		defer func() {
			if err := shutdown(); err != nil {
				slog.Error("Error shutting down OpenTelemetry", "error", err)
			}
		}()
	}

	// Connection string for PostgreSQL
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	// Connect to the PostgreSQL database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		panic(err)
	}
	slog.Info("Connected to PostgreSQL database", "host", cfg.Database.Host, "database", cfg.Database.Name)

	// Connect to the Redis container
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}
	slog.Info("Connected to Redis", "addr", cfg.Redis.Addr)

	// Init Snowflake node
	n, err := snowflake.NewNode(cfg.SnowFlake.Node)
	if err != nil {
		slog.Error("Failed to create a snowflake node", "error", err)
		panic(err)
	}
	slog.Info("Snowflake node initialized", "node", cfg.SnowFlake.Node)

	s := storage.New(db)
	svc := service.New(s, rdb)
	h := handler.New(svc, n)

	r := NewRouter(h)
	slog.Info("Server starting", "port", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		slog.Error("Server failed to start", "error", err)
		panic(err)
	}
}
