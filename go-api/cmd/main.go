// Package main es el punto de entrada de la API Go.
//
// @title           QR Matrix API
// @version         1.0
// @description     API REST que calcula la factorización QR de una matriz y obtiene estadísticas desde Node.js.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.email  support@example.com
//
// @license.name  MIT
//
// @host      localhost:8080
// @BasePath  /
package main

import (
	"log"

	_ "go-api/docs"
	"go-api/internal/clients"
	"go-api/internal/config"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"go-api/internal/routes"
	"go-api/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/gofiber/swagger"
)

func main() {
	// Cargar configuración desde .env / variables de entorno.
	cfg := config.Load()

	// Componer dependencias (dependency injection manual).
	nodeClient := clients.NewNodeClient(cfg.NodeAPIURL)
	matrixSvc := services.NewMatrixService(nodeClient)
	matrixHandler := handlers.NewMatrixHandler(matrixSvc)

	// Inicializar Fiber con el manejador centralizado de errores.
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middlewares globales.
	app.Use(middleware.Recovery())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New())

	// Swagger UI disponible en /swagger/index.html
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// Health-check mínimo.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "env": cfg.AppEnv})
	})

	// Registrar rutas de la API.
	routes.Register(app, matrixHandler)

	log.Printf("[main] starting server on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("[main] server error: %v", err)
	}
}
