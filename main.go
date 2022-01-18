package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gastonstec/utils/config"

	"gastonstec/nuricc/routes"
)

func main () {
	var err error

	// Load config file values
	err = config.LoadConfig("/", false)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	// Create fiber application object
	app := fiber.New()
	app.Use(cors.New()) // enable CORS middleware

	// Catch shutdown application signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Shutdown application started...")
		err = app.Shutdown()
		log.Println("application ended")
	}()

	routes.InitRoutes(app)

	// Start fiber listening
	err = app.Listen(config.GetString("APP_PORT"))

	if err != nil {
		log.Println(err.Error())
	}

	// Cleanup tasks
	log.Println("Application cleanup tasks")

}