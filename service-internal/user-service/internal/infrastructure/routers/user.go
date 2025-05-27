package routers

import (
	"github.com/khuongdo95/go-service/internal/adapter/controllers"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
)

func UserRouters(app *global.HttpServer, userController *controllers.UserController) {
	app.Get("/users/:id", userController.Get)
	app.Get("/users", userController.List)
	app.Post("/users", userController.Create)
	app.Put("/users/:id", userController.Update)
	app.Delete("/users/:id", userController.Delete)
	app.Get("/me", userController.Me)
}
