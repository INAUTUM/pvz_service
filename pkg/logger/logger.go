// pkg/logger/logger.go
package logger

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	logger, _ := zap.NewProduction()
	return &Logger{logger}
}

func (l *Logger) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Реализация middleware
			return next(c)
		}
	}
}

// func New() *zap.Logger {
// 	config := zap.NewProductionConfig()
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	logger, _ := config.Build()
// 	return logger
// }

func Middleware(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			latency := time.Since(start)
			log.Info("Request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", latency),
			)
			
			return err
		}
	}
}

func GRPCLoggerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		log.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("latency", time.Since(start)),
			zap.Error(err),
		)
		return
	}
}

func (l *Logger) MiddlewareV4() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Реализация middleware для echo v4
			start := time.Now()
			err := next(c)
			l.Info("Request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", c.Response().Status),
				zap.Duration("duration", time.Since(start)),
			)
			return err
		}
	}
}