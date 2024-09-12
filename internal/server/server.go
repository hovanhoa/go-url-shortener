package server

import (
	"fmt"
	"github.com/hovanhoa/go-url-shortener/config"
	"github.com/hovanhoa/go-url-shortener/internal/handler"
	"github.com/hovanhoa/go-url-shortener/internal/service"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
	"github.com/hovanhoa/go-url-shortener/pkg/snowflake"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"

	_ "github.com/lib/pq" // Import the pq driver
)

func Init() {
	cfg := config.GetConfig()
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
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	n, err := snowflake.NewNode(cfg.SnowFlake.Node)
	if err != nil {
		log.Fatal("Failed to create a snowflake node:", err)
	}

	s := storage.New(db)
	svc := service.New(s)
	h := handler.New(svc, n)

	r := NewRouter(h)
	r.Run(cfg.Server.Port)
}
