package core

import (
	"strings"
	"time"

	"github.com/khuongdo95/go-pkg/appLogger"
	"github.com/khuongdo95/go-pkg/common/response"

	"github.com/khuongdo95/go-service/internal/adapter/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"

	"github.com/khuongdo95/go-pkg/fiberzap-logger"

	"github.com/gofiber/fiber/v3/middleware/recover"
)

// TODO this core will not in package
func NewHttpServer(allowOrigins string, log *appLogger.Logger) (*fiber.App, *response.AppError) {
	// Configure Fiber app

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.Use(
		helmet.New(
			helmet.Config{
				CrossOriginEmbedderPolicy: "false",
				ContentSecurityPolicy:     "false",
				CrossOriginOpenerPolicy:   "cross-origin",
				CrossOriginResourcePolicy: "cross-origin",
			},
		),
		cors.New(
			cors.Config{
				AllowOrigins:     strings.Split(allowOrigins, ","),
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				AllowCredentials: false,
			},
		),
		recover.New(),
		// update new version
		fiberzap.New(
			fiberzap.Config{
				Logger: log.Logger(),
			},
		),
		middlewares.ReqContextHandler,
		middlewares.LoggingInterceptor(log),
	)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("I'm good!")
	})

	// app.Use(func() fiber.Handler {
	// 	return func(c fiber.Ctx) error {
	// 		return c.Status(http.StatusNotFound).JSON(common.ErrorResponse(helper.GetHttpReqCtx(c), response.NotFound("route not found")))
	// 	}
	// }())

	return app, nil
}
