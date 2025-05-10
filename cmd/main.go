package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"github/ramabmtr/billing-engine/config"
	_ "github/ramabmtr/billing-engine/docs"
	"github/ramabmtr/billing-engine/internal/handler"
	"github/ramabmtr/billing-engine/internal/repository"
	"github/ramabmtr/billing-engine/internal/service"
)

// @title Billing Engine API
// @version 1.0

// @contact.name Rama Bramantara
// @contact.email ramabmtr@gmail.com

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	// init configuration
	config.InitEnv()
	config.InitDB()

	// Initialize repositories
	borrowerRepo := repository.NewBorrowerRepo(config.GetDB())
	loanRepo := repository.NewLoanRepo(config.GetDB())
	loanPaymentRepo := repository.NewLoanPaymentRepo(config.GetDB())

	// Initialize services
	borrowerSvc := service.NewBorrowerService(borrowerRepo)
	loanSvc := service.NewLoanService(loanRepo, loanPaymentRepo)

	// Initialize handlers
	borrowerHandler := handler.NewBorrowerHandler(borrowerSvc)
	loanHandler := handler.NewLoanHandler(loanSvc)
	paymentHandler := handler.NewPaymentHandler(loanSvc)

	// Initialize Echo
	e := echo.New()
	e.Validator = config.NewValidator()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-KEY",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == config.GetEnv().Server.ApiKey, nil
		},
		Skipper: func(c echo.Context) bool {
			skippedPaths := []string{"/api/ping", "/docs"}
			for _, path := range skippedPaths {
				if strings.HasPrefix(c.Request().URL.Path, path) {
					return true
				}
			}
			return false
		},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			skippedPaths := []string{"/docs"}
			for _, path := range skippedPaths {
				if strings.HasPrefix(c.Request().URL.Path, path) {
					return true
				}
			}
			return false
		},
		Format:           middleware.DefaultLoggerConfig.Format,
		CustomTimeFormat: middleware.DefaultLoggerConfig.CustomTimeFormat,
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API routes
	e.GET("/docs/*", echoSwagger.WrapHandler)

	apiGroup := e.Group("/api")
	// @Summary Health check endpoint
	// @Description Returns a simple pong message to verify the API is running
	// @Tags system
	// @Produce json
	// @Success 200 {object} map[string]string "Returns message: pong"
	// @Router /ping [get]
	apiGroup.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})
	borrowerHandler.RegisterRoutes(apiGroup)
	loanHandler.RegisterRoutes(apiGroup)
	paymentHandler.RegisterRoutes(apiGroup)

	// Start server
	serverAddr := fmt.Sprintf(":%d", config.GetEnv().Server.Port)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
