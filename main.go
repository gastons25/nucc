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

	// Open database connection and defer close
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

	// Initialize fiber application routes
	routes.InitRoutes(app)

	// Catch shutdown application signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Application shutting down")
		err = app.Shutdown()
	}()

	// Start fiber listening
	err = app.Listen(config.GetString("APP_PORT"))
	if err != nil {
		log.Println(err.Error())
	}

	// Final cleanup tasks

}