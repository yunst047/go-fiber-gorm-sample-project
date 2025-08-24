package main

import (
	"log"
	"os"

	"go-fiber-gorm-sample/config"
	"go-fiber-gorm-sample/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	db, err := database.NewMySQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	app := fiber.New()

	app.Use(logger.New())

	//api := app.Group("/api")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = config.Port
	}

	log.Printf("ASVL Ticket API Server started on http://0.0.0.0:%s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
