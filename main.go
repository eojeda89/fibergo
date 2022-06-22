package main

import (
	"fibergo/configs"
	"fibergo/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{"data": "Hello from fiber & mongo"})
	})

	configs.ConnectDB()
	routes.UserRoute(app)
	app.Listen("https://sleepy-harbor-61821.herokuapp.com:8443")
}
