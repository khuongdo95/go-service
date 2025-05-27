package middlewares

import (
	"reflect"

	"github.com/khuongdo95/go-service/internal/infrastructure/helper"

	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"

	"github.com/gofiber/fiber/v3"
)

// TODO if Error status 500 to send notfi discord
func ErrorHandler(c fiber.Ctx, err error) error {
	if err != nil && reflect.ValueOf(err).Kind() == reflect.Pointer {
		internalError := response.NewAppError(err.Error(), constants.InternalServerErr)
		if appErr, ok := err.(*response.AppError); ok {
			internalError = appErr
		}

		return c.Status(constants.HttpCode[internalError.Code]).JSON(common.ErrorResponse(helper.GetHttpReqCtx(c), internalError))
	}

	return nil
}
