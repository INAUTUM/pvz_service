package http

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"pvz_service/internal/domain"
	"pvz_service/internal/service"
	"pvz_service/pkg/jwt"
)

type Handler struct {
	svc          Service
	jwtSecret    string
	metrics      *service.Metrics
}

type Service interface {
	CreatePVZ(ctx context.Context, pvz domain.PVZ) (*domain.PVZ, error)
	StartReception(ctx context.Context, pvzID string) (*domain.Reception, error)
	AddProduct(ctx context.Context, receptionID string, product domain.Product) error
	CloseReception(ctx context.Context, pvzID string) error
	GetPVZs(ctx context.Context, filter domain.PVZFilter) ([]domain.PVZ, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error
}

func NewHandler(svc Service, jwtSecret string, metrics *service.Metrics) *Handler {
	return &Handler{
		svc:       svc,
		jwtSecret: jwtSecret,
		metrics:   metrics,
	}
}

func (h *Handler) PostDummyLogin(c echo.Context) error {
	var req PostDummyLoginJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "invalid request"})
	}

	token, err := jwt.GenerateToken(string(req.Role), h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "failed to generate token"})
	}

	return c.JSON(http.StatusOK, Token(token))
}

func (h *Handler) PostRegister(c echo.Context) error {
	var req PostRegisterJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "invalid request"})
	}

	if !domain.UserRole(req.Role).IsValid() {
		return c.JSON(http.StatusBadRequest, Error{Message: "invalid user role"})
	}

	// Реализация регистрации пользователя
	// ...

	return c.JSON(http.StatusCreated, User{
		Id:    &userID,
		Email: req.Email,
		Role:  domain.UserRole(req.Role),
	})
}

func (h *Handler) PostPvz(c echo.Context) error {
	claims, err := h.getClaims(c)
	if err != nil || claims.Role != "moderator" {
		return c.JSON(http.StatusForbidden, Error{Message: "access denied"})
	}

	var pvz PVZ
	if err := c.Bind(&pvz); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "invalid request"})
	}

	createdPVZ, err := h.svc.CreatePVZ(c.Request().Context(), domain.PVZ{
		City: pvz.City,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	h.metrics.PVZCreated.Inc()
	return c.JSON(http.StatusCreated, createdPVZ)
}

func (h *Handler) PostReceptions(c echo.Context) error {
    claims, err := h.getClaims(c)
    if err != nil || claims.Role != "employee" {
        return c.JSON(http.StatusForbidden, Error{Message: "access denied"})
    }

    var req PostReceptionsJSONRequestBody
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, Error{Message: "invalid request"})
    }

    // Преобразование и валидация UUID
    pvzID, err := uuid.Parse(req.PvzId.String())
    if err != nil {
        return c.JSON(http.StatusBadRequest, Error{Message: "invalid pvz id format"})
    }

    // Использование pvzID в вызове сервиса
    reception, err := h.svc.StartReception(
        c.Request().Context(), 
        pvzID.String(), // Добавляем преобразование в строку
    )
    if err != nil {
        return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
    }

    h.metrics.ReceptionsCreated.Inc()
    return c.JSON(http.StatusCreated, reception)
}

func (h *Handler) PostProducts(c echo.Context) error {
	claims, err := h.getClaims(c)
	if err != nil || claims.Role != "employee" {
		return c.JSON(http.StatusForbidden, Error{Message: "access denied"})
	}

	var req PostProductsJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "invalid request"})
	}

	productType, err := domain.ParseProductType(req.Type)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	err = h.svc.AddProduct(c.Request().Context(), req.PvzId.String(), domain.Product{
		Type:      productType, 
		AddedAt:   time.Now(),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	h.metrics.ProductsAdded.Inc()
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) getClaims(c echo.Context) (*jwt.Claims, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	return jwt.ValidateToken(token, h.jwtSecret)
}

func MetricsMiddleware(metrics *service.Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Реализация сбора метрик
			return next(c)
		}
	}
}

func RegisterHandlers(e *echo.Echo, h *Handler) {
	// Регистрация всех маршрутов
	e.POST("/dummyLogin", h.PostDummyLogin)
	e.POST("/register", h.PostRegister)
	e.POST("/pvz", h.PostPvz)
	// Остальные маршруты...
}