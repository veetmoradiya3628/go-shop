package main

import (
	"github.com/gin-gonic/gin"
	"github.com/veetmoradiya3628/go-shop/internal/config"
	"github.com/veetmoradiya3628/go-shop/internal/database"
	"github.com/veetmoradiya3628/go-shop/internal/logger"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database instance")
	}
	defer mainDB.Close()
	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("starting server")
}
