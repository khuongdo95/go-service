package helper

import (
	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/constants"

	"github.com/gofiber/fiber/v3"
)

func GetHttpUserCtx(c fiber.Ctx) *common.UserInfo[any] {
	return c.Locals(constants.REQUEST_CONTEXT_KEY).(*common.ReqContext).UserInfo
}

func GetHttpReqCtx(c fiber.Ctx) *common.ReqContext {
	return c.Locals(constants.REQUEST_CONTEXT_KEY).(*common.ReqContext)
}

func SetHttpUserCtx(c fiber.Ctx, user *common.UserInfo[any]) {
	c.Locals(constants.REQUEST_CONTEXT_KEY).(*common.ReqContext).UserInfo = user
}
