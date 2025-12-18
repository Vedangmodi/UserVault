package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"uservault/config"
	"uservault/internal/handler"
	applogger "uservault/internal/logger"
	"uservault/internal/middleware"
	"uservault/internal/repository"
	"uservault/internal/routes"
	"uservault/internal/service"
)

func main() {
	cfg := config.Load()

	zapLogger, err := applogger.NewLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer zapLogger.Sync()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("failed to create pgx pool", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		zapLogger.Fatal("failed to ping database", zap.Error(err))
	}

	app := fiber.New()

	// Middlewares
	app.Use(middleware.ZapLogger(zapLogger))

	// Dependencies
	repo := repository.NewUserRepository(pool)
	validate := validator.New()
	svc := service.NewUserService(repo, validate)
	userHandler := handler.NewUserHandler(svc, zapLogger)

	// Routes
	routes.Register(app, userHandler)

	zapLogger.Sugar().Infof("Starting server on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		zapLogger.Sugar().Fatalf("failed to start server: %v", err)
	}
}


