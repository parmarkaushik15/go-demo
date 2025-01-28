package router

import (
	"go-demo-api/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/user")
	// routes
	v1.Get("/", handler.GetAllUsers) 
}