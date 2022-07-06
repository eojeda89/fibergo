package main

import (
	"fibergo/configs"
	_ "fibergo/docs"
	"fibergo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"log"
	"os"
)

// @title Fiber User Microservice API
// @version 1.0
// @description This is a microservice for the user´s management
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	port := os.Getenv("PORT")
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{"data": "Hello from fiber & mongo"})
	})

	configs.ConnectDB()
	routes.UserRoute(app)
	port = "8080" //dev only
	log.Fatal(app.Listen(":" + port))
}
