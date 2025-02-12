package main

import (
	"coffee-shop-api/config"
	"coffee-shop-api/database"
	"coffee-shop-api/internal/app"
	"coffee-shop-api/monitoring"
	"context"
	"log"
)

// @title Coffee Shop
// @version 1.0
// @description This is coffee shop apis documentations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email fiber@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
func main() {
	// Init
	cfg := config.MustLoadConfig()
	logger := monitoring.NewLogger(cfg)
	defer logger.GetLogger().Sync()

	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := monitoring.SetupOTelSDK(ctx, cfg)
	if err != nil {
		return
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			logger.Error(ctx, "failed to shutdown OpenTelemetry", err)
		}
	}()

	db := database.MustNewDatabase(cfg, logger)
	defer db.Close()

	application := app.New(cfg, db.DB(), logger)
	application.RegisterHandlers()
	log.Fatal(application.Start())
}
