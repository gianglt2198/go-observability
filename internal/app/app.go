package app

import (
	"coffee-shop-api/config"
	"coffee-shop-api/internal/handlers"
	"coffee-shop-api/middlewares"
	"coffee-shop-api/monitoring"
	"coffee-shop-api/routes"
	"fmt"

	_ "coffee-shop-api/docs"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

type App struct {
	cfg      *config.Config
	db       *gorm.DB
	app      *fiber.App
	logger   *monitoring.AppLogger
	handlers []Handler
}

type Handler interface {
	Register(router fiber.Router)
}

func New(cfg *config.Config, db *gorm.DB, logger *monitoring.AppLogger) *App {

	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middlewares.RequestIDMiddleware)
	app.Use(middlewares.TracingMiddleware("main", "request_caller",
		middlewares.TracingConfig{
			ServiceName:    cfg.App.Name,
			ServiceVersion: "1.0.0",
		}))
	app.Use(middlewares.MetricMiddleware(middlewares.MetricConfig{
		ServiceName:    cfg.App.Name,
		ServiceVersion: "1.0.0",
	}))

	prometheus := fiberprometheus.New(cfg.App.Name)
	prometheus.RegisterAt(app, "/metrics")
	prometheus.SetSkipPaths([]string{
		"/metrics", "/health", "/swagger",
	})

	// app.Use(prometheus.Middleware)

	app.Get("/health", HealthCheck)

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8081/swagger/oauth2-redirect.html",
	}))

	return &App{
		cfg:    cfg,
		db:     db,
		app:    app,
		logger: logger,
	}
}

func (a *App) RegisterHandlers() {
	a.handlers = []Handler{
		handlers.NewCoffeeHandler(a.cfg, a.db, a.logger),
	}

	api := a.app.Group("/api")
	for _, h := range a.handlers {
		h.Register(api)
	}

	a.app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(
			routes.ErrorResponse(fiber.ErrBadGateway),
		)
	})
}

func (a *App) Start() error {
	return a.app.Listen(fmt.Sprintf(":%d", a.cfg.App.Port))
}

func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
