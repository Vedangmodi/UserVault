package routes

import (
	"github.com/gofiber/fiber/v2"

	"uservault/internal/handler"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/:id", userHandler.GetUser)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
	app.Get("/users", userHandler.ListUsers)
}


