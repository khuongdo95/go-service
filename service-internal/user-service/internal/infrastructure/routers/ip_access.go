package routers

import (
	"github.com/khuongdo95/go-service/internal/adapter/controllers"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
)

func IpAccessRouters(app *global.HttpServer, userController *controllers.IpAccessController) {
	app.Post("/ip", userController.CreateIP)
	app.Get("/ip/:id", userController.GetIP)
	app.Post("/ip/list", userController.ListIP)
	app.Put("/ip/:id", userController.UpdateIP)
	app.Delete("/ip/:id", userController.DeleteIP)
}
