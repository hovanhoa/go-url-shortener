package server

import (
	"database/sql"
	"fmt"
	"github.com/hovanhoa/go-url-shortener/config"
	"github.com/hovanhoa/go-url-shortener/internal/handler"
	"github.com/hovanhoa/go-url-shortener/internal/service"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
	"log"

	_ "github.com/lib/pq" // Import the pq driver
)

func Init() {
	cfg := config.GetConfig()
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	// Open a connection to the database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error on connecting to database, %v", err)
	}

	// Ensure the connection is closed when the function exits
	defer conn.Close()

	s := storage.New(conn)
	svc := service.New(s)
	h := handler.New(svc)

	r := NewRouter(h)
	r.Run(cfg.Server.Port)
}
