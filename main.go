package main

import (
	"log"

	"go-fiber-gorm-sample/config"
	"go-fiber-gorm-sample/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {

	db, err := database.NewMySQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/metrics", monitor.New())

	//api := app.Group("/api")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := config.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Sample API Server started on http://0.0.0.0:%s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
