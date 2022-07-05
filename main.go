package main

import (
	"fibergo/configs"
	"fibergo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{"data": "Hello from fiber & mongo"})
	})

	configs.ConnectDB()
	routes.UserRoute(app)
	port = "8080" //dev only
	log.Fatal(app.Listen(":" + port))
}
