package main

import (
	"fibergo/configs"
	"fibergo/routes"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{"data": "Hello from fiber & mongo"})
	})

	configs.ConnectDB()
	routes.UserRoute(app)
	log.Fatal(app.Listen(":" + port))
}
