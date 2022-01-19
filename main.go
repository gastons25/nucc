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

	log.Println(utils.GetOsEnv())
	log.Println(utils.GetGolangEnv())

	// Load config file values
	err = config.LoadConfig("/", false)
	if err != nil {
		log.Println(utils.GetFunctionName() + ": " + err.Error())
		os.Exit(1)
	}

	// Open database connection
	err = db.OpenDB()
	if err != nil {
		log.Println(utils.GetFunctionName() + ": " + err.Error())
		os.Exit(1)
	}
	defer db.CloseDB()

	// Load memdb
	err = memdb.Load()
	if err != nil {
		log.Println(utils.GetFunctionName() + ": " + err.Error())
		os.Exit(1)
	}
	
	// Create fiber application object
	app := fiber.New()
	// enable CORS middleware
	app.Use(cors.New())

	// Catch shutdown application signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Application shutting down")
		err = app.Shutdown()
	}()

	// Initialize routes
	routes.InitRoutes(app)

	// Start fiber listening
	err = app.Listen(config.GetString("APP_PORT"))
	if err != nil {
		log.Println(err.Error())
	}

	// Cleanup tasks

}