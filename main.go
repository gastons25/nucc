// Package main implements the main point of entry.
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gastonstec/utils/config"
	"github.com/gastonstec/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"


	"gastonstec/nuricc/db"
	"gastonstec/nuricc/memdb"
	"gastonstec/nuricc/routes"

)

// Function main handles the initial settings and start
func main () {
	var err error

	// Log environment variables
	log.Println(utils.GetOsEnv())
	log.Println(utils.GetGolangEnv())

	// Load config file values from .env file
	err = config.LoadConfig("/", false)
	if err != nil {
		log.Println(utils.GetFunctionName() + ": " + err.Error())
		os.Exit(1)
	}

	// Open database connection pool and defer close
	err = db.OpenDB()
	if err != nil {
		log.Println(utils.GetFunctionName() + ": " + err.Error())
		os.Exit(1)
	}
	defer db.CloseDB()

	// Load in-memory database
	err = memdb.Load()
	if err != nil {
		log.Println(utils.GetFunctionName() + ": " + err.Error())
		os.Exit(1)
	}
	
	// Create fiber application object
	app := fiber.New()
	app.Use(cors.New()) // enable cors

	// Initialize fiber routes
	routes.InitRoutes(app)

	// Catch shutdown application signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Shutting down application")
		err = app.Shutdown()
	}()

	// Start fiber listening
	err = app.Listen(config.GetString("APP_PORT"))
	if err != nil {
		log.Println(err.Error())
	}

	// Final cleanup tasks
	log.Println("Final tasks")

}