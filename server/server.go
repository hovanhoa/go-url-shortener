package server

import "github.com/hovanhoa/go-url-shortener/config"

func Init() {
	cfg := config.GetConfig()
	r := NewRouter()
	r.Run(cfg.Server.Port)
}
