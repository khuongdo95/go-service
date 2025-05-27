package middlewares

import (
	"github.com/khuongdo95/go-pkg/appLogger"
	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"

	"github.com/gofiber/fiber/v3"
)

func LoggingInterceptor(log *appLogger.Logger) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		requestContext := c.Locals(constants.REQUEST_CONTEXT_KEY).(*common.ReqContext)
		log.ReqClientLog(requestContext, c.Method(), c.Path())

		err := c.Next()
		errApp := response.ConvertError(err)
		if errApp != nil {
			log.ResClientLog(requestContext, uint(c.Response().StatusCode()), errApp)
		} else {
			log.ResClientLog(requestContext, uint(c.Response().StatusCode()), err)
		}

		return err
	}
}
