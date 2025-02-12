package handlers

import (
	"coffee-shop-api/config"
	"coffee-shop-api/internal/common"
	"coffee-shop-api/internal/dto"
	"coffee-shop-api/internal/models"
	"coffee-shop-api/middlewares"
	"coffee-shop-api/monitoring"
	"coffee-shop-api/routes"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	CoffeeHandler = "coffee_handler"
)

type coffeeHandler struct {
	cfg    *config.Config
	db     *gorm.DB
	logger *monitoring.AppLogger
}

func NewCoffeeHandler(cfg *config.Config,
	db *gorm.DB,
	logger *monitoring.AppLogger,
) *coffeeHandler {
	return &coffeeHandler{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func (h *coffeeHandler) Register(router fiber.Router) {
	coffee := router.Group("/coffee")

	coffee.Post("/", routes.Usecase(h.CreateCoffee,
		fiber.StatusCreated,
		middlewares.AllPayloadValidator[dto.CreateCoffeeRequest](),
	))
	coffee.Get("/", routes.Usecase(h.ListCoffee,
		fiber.StatusOK,
		middlewares.AllPayloadValidator[dto.ListCoffeeRequest](),
	))
	coffee.Get("/:id", routes.Usecase(h.GetCoffee,
		fiber.StatusOK,
		middlewares.AllPayloadValidator[dto.GetCoffeeRequest](),
	))
	coffee.Patch("/:id", routes.Usecase(h.UpdateCoffee,
		fiber.StatusOK,
		middlewares.AllPayloadValidator[dto.UpdateCoffeeRequest](),
	))
	coffee.Delete("/:id", routes.Usecase(h.DeleteCoffee,
		fiber.StatusOK,
		middlewares.AllPayloadValidator[dto.DeleteCoffeeRequest](),
	))
}

// CreateCoffee godoc
// @Router             /api/coffee [post]
// @Summary            CreateCoffee
// @Description        CreateCoffee
// @Tags               coffee
// @Accept             json
// @Produce            json
// @Param              request     body          dto.CreateCoffeeRequest          true    "request body"
// @Success            200         {object}      models.Coffee
// @Failure            400         {object}      error
// @Failure            404         {object}      error
// @Failure            500         {object}      error
func (h *coffeeHandler) CreateCoffee(ctx context.Context, req dto.CreateCoffeeRequest) (*models.Coffee, error) {
	ctx, span := middlewares.GetSpanFromContext(ctx, "create_coffee")
	defer span.End()

	h.logger.Info(ctx, "[coffeeHandler]CreateCoffee", req)

	entity := &models.Coffee{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
	if err := h.db.WithContext(ctx).Create(entity).Error; err != nil {
		h.logger.Error(ctx, "[coffeeHandler]CreateCoffee: failed to create coffee", err)
		return nil, err
	}
	return entity, nil
}

// UpdateCoffee godoc
// @Router              /api/coffee/{id} [patch]
// @Summary             UpdateCoffee
// @Description         UpdateCoffee
// @Tags                coffee
// @Accept              json
// @Produce             json
// @Param               id              path            string                          true    "id"
// @Param               request         body            dto.UpdateCoffeeRequest         true    "request body"
// @Success             200             {object}        models.Coffee
// @Failure             400             {object}        error
// @Failure             404             {object}        error
// @Failure             500             {object}        error
func (h *coffeeHandler) UpdateCoffee(ctx context.Context, req dto.UpdateCoffeeRequest) (*models.Coffee, error) {
	ctx, span := middlewares.GetSpanFromContext(ctx, "update_coffee")
	defer span.End()

	entity := &models.Coffee{}

	if req.Description != "" {
		entity.Description = req.Description
	}

	if req.Name != "" {
		entity.Name = req.Name
	}

	if req.Price > 0 {
		entity.Price = req.Price
	}

	if err := h.db.WithContext(ctx).Model(entity).Where("id = ?", req.ID).Updates(entity).Error; err != nil {
		h.logger.Error(ctx, "[coffeeHandler]UpdateCoffee: failed to update coffee", err)
		return nil, err
	}
	return entity, nil
}

// DeleteCoffee godoc
// @Router          /api/coffee/{id} [delete]
// @Summary         DeleteCoffee
// @Description     DeleteCoffee
// @Tags            coffee
// @Accept          json
// @Produce         json
// @Param           id              path            string              true        "id"
// @Success         200             {object}        bool
// @Failure         400             {object}        error
// @Failure         404             {object}        error
// @Failure         500             {object}        error
func (h *coffeeHandler) DeleteCoffee(ctx context.Context, req dto.DeleteCoffeeRequest) (*bool, error) {
	ctx, span := middlewares.GetSpanFromContext(ctx, "delete_coffee")
	defer span.End()

	if err := h.db.WithContext(ctx).Delete("id = ?", req.ID).Error; err != nil {
		h.logger.Error(ctx, "[coffeeHandler]DeleteCoffee: failed to delete coffee", err)
		return nil, err
	}
	return lo.ToPtr(true), nil
}

// GetCoffee godoc
// @Router          /api/coffee/{id} [get]
// @Summary         GetCoffee
// @Description     GetCoffee
// @Tags            coffee
// @Accept          json
// @Produce         json
// @Param           id              path            string              true        "id"
// @Success         200             {object}        models.Coffee
// @Failure         400             {object}        error
// @Failure         404             {object}        error
// @Failure         500             {object}        error
func (h *coffeeHandler) GetCoffee(ctx context.Context, req dto.GetCoffeeRequest) (*models.Coffee, error) {
	ctx, span := middlewares.GetSpanFromContext(ctx, "get_coffee")
	defer span.End()

	h.logger.Info(ctx, "[coffeeHandler]GetCoffee", req)
	var entity models.Coffee
	if err := h.db.WithContext(ctx).Model(&entity).Where("id = ?", req.ID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewError(fiber.ErrNotFound, err)
		}
		h.logger.Error(ctx, "[coffeeHandler]GetCoffee: failed to get coffee", err)
		return nil, common.NewError(fiber.ErrInternalServerError, err)
	}
	return lo.ToPtr(entity), nil
}

// ListCoffee godoc
// @Router              /api/coffee [get]
// @Summary             ListCoffee
// @Description         ListCoffee
// @Tags                coffee
// @Accept              json
// @Produce             json
// @Param               payload         query           dto.ListCoffeeRequest       true     "request payload"
// @Success             200             {object}        models.Coffee
// @Failure             400             {object}        error
// @Failure             404             {object}        error
// @Failure             500             {object}        error
func (h *coffeeHandler) ListCoffee(ctx context.Context, req dto.ListCoffeeRequest) ([]*models.Coffee, error) {
	// Create a new context with parent span
	ctx, span := middlewares.GetSpanFromContext(ctx, "list_coffee")
	defer span.End()

	h.logger.Info(ctx, "[coffeeHandler]ListCoffee", req)

	var entities []*models.Coffee
	if err := h.db.WithContext(ctx).Where("1=1").Find(&entities).Error; err != nil {
		h.logger.Error(ctx, "[coffeeHandler]ListCoffee: failed to list coffee", err)
		return nil, err
	}

	return entities, nil
}
