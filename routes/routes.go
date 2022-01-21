// Package routes implements the routing settings.
package routes

import (
	"gastonstec/nuricc/handlers"
	"github.com/gofiber/fiber/v2"
)

// Function InitRoutes set the endpoints routes and handlers
func InitRoutes(app *fiber.App) {

	// Get block route
	app.Get("/api/v1/block/:network_code/:block_hash", handlers.GetBlock)
	// Get transaction route
	app.Get("/api/v1/tx/:network_code/:transaction_id", handlers.GetTx)
	
}