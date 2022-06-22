package routes

import (
	"fibergo/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Get("/", controllers.HelloWorld)
	app.Post("/users", controllers.CreateUser)
	app.Post("/users/flat", controllers.CreateUser2)
	app.Get("/users/:userId", controllers.GetAUser)
	app.Get("/users", controllers.GetAllUsers)
	app.Put("/users/:userId", controllers.EditAUser)
	app.Delete("/users/:userId", controllers.DeleteAUser)
}
