// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpcCtrl "pvz_service/internal/controller/grpc"
	httpCtrl "pvz_service/internal/controller/http"
	"pvz_service/internal/domain"
	"pvz_service/internal/service"
	"pvz_service/internal/storage/postgres"
	"pvz_service/pkg/logger"
)

type Application struct {
	cfg        *Config
	echo       *echo.Echo
	grpcServer *grpc.Server
	db         *sqlx.DB
	metrics    *service.Metrics
	logger     *logger.Logger
}

func NewApplication(cfg *Config) *Application {
	return &Application{
		cfg:    cfg,
		logger: logger.New(),
	}
}

func (a *Application) initDB() error {
	db, err := sqlx.Connect("postgres", a.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	a.db = db
	return nil
}

func (a *Application) initRepositories() (
	domain.UserRepository,
	domain.PVZRepository,
	domain.ReceptionRepository,
	domain.ProductRepository,
) {
	return postgres.NewUserRepository(a.db),
		postgres.NewPVZRepository(a.db),
		postgres.NewReceptionRepository(a.db),
		postgres.NewProductRepository(a.db)
}

func (a *Application) initServices(
	userRepo domain.UserRepository,
	pvzRepo domain.PVZRepository,
	receptionRepo domain.ReceptionRepository,
	productRepo domain.ProductRepository,
) (service.AuthService, *service.PVZService, *service.ReceptionService) {
	
	metrics := service.NewMetrics()
	a.metrics = metrics

	return service.NewAuthService(userRepo, a.cfg.JWTSecret),
		service.NewPVZService(pvzRepo),
		service.NewReceptionService(receptionRepo, productRepo)
}

func (a *Application) initHTTPServer(
	authService service.AuthService,
	pvzService *service.PVZService,
	receptionService *service.ReceptionService,
) {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(a.logger.MiddlewareV4())
	e.Use(httpCtrl.MetricsMiddleware(a.metrics))

	serviceAdapter := &serviceAdapter{
		auth:       authService,
		pvz:        pvzService,
		reception:  receptionService,
	}

	handler := httpCtrl.NewHandler(
		serviceAdapter,
		a.cfg.JWTSecret,
		a.metrics,
	)

	httpCtrl.RegisterHandlers(e, handler)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	a.echo = e
}

func (a *Application) initGRPCServer(pvzRepo domain.PVZRepository) {
	grpcController := grpcCtrl.NewPVZController(pvzRepo)
	a.grpcServer = grpc.NewServer()
	grpcController.Register(a.grpcServer)
}

func (a *Application) Run() error {
	if err := a.initDB(); err != nil {
		return err
	}
	defer a.db.Close()

	if err := postgres.RunMigrations(a.db.DB); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	userRepo, pvzRepo, receptionRepo, productRepo := a.initRepositories()

	authService, pvzService, receptionService := a.initServices(
		userRepo,
		pvzRepo,
		receptionRepo,
		productRepo,
	)

	a.initHTTPServer(authService, pvzService, receptionService)
	a.initGRPCServer(pvzRepo)

	ctx, stop := signal.NotifyContext(context.Background(), 
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		a.logger.Info("Starting HTTP server", 
			zap.String("port", a.cfg.HTTPPort))
		if err := a.echo.Start(":" + a.cfg.HTTPPort); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("HTTP server error", 
				zap.Error(err))
		}
	}()

	go func() {
		a.logger.Info("Starting gRPC server", 
			zap.String("port", a.cfg.GRPCPort))
		listener, err := net.Listen("tcp", ":"+a.cfg.GRPCPort)
		if err != nil {
			a.logger.Fatal("gRPC listener error", 
				zap.Error(err))
		}
		
		if err := a.grpcServer.Serve(listener); err != nil {
			a.logger.Fatal("gRPC server error", 
				zap.Error(err))
		}
	}()

	<-ctx.Done()
	a.logger.Info("Shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.echo.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", 
			zap.Error(err))
	}

	a.grpcServer.GracefulStop()
	a.logger.Info("Servers stopped properly")
	return nil
}

// Полная реализация интерфейса Service
type serviceAdapter struct {
	auth      service.AuthService
	pvz       *service.PVZService
	reception *service.ReceptionService
}

func (a *serviceAdapter) CreatePVZ(ctx context.Context, pvz domain.PVZ) (*domain.PVZ, error) {
	return a.pvz.CreatePVZ(ctx, pvz)
}

func (a *serviceAdapter) GetPVZs(ctx context.Context, filter domain.PVZFilter) ([]domain.PVZ, error) {
	return a.pvz.GetPVZs(ctx, filter)
}

func (a *serviceAdapter) AddProduct(ctx context.Context, receptionID string, product domain.Product) error {
	uuid, err := uuid.Parse(receptionID)
	if err != nil {
		return domain.ErrInvalidID
	}
	return a.reception.AddProduct(ctx, uuid, product)
}

func (a *serviceAdapter) CloseReception(ctx context.Context, pvzID string) error {
	uuid, err := uuid.Parse(pvzID)
	if err != nil {
		return domain.ErrInvalidID
	}
	return a.reception.CloseReception(ctx, uuid)
}

// Реализация остальных методов интерфейса
func (a *serviceAdapter) StartReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	uuid, err := uuid.Parse(pvzID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}
	return a.reception.StartReception(ctx, uuid)
}

func (a *serviceAdapter) DeleteLastProduct(ctx context.Context, pvzID string) error {
    uuid, err := uuid.Parse(pvzID)
    if err != nil {
        return domain.ErrInvalidID
    }
    
    // Получаем активную приемку
    reception, err := a.reception.GetActiveReception(ctx, uuid)
    if err != nil {
        return err
    }
    
    return a.reception.DeleteLastProduct(ctx, reception.ID)
}