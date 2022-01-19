package routes

import (
	"gastonstec/nuricc/handlers"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {

	app.Get("/api/v1/get_block/:network_code/:block_hash", handlers.GetBlock)
	
}