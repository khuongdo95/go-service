package middlewares

import (
	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/constants"

	"github.com/gofiber/fiber/v3"
)

func ReqContextHandler(c fiber.Ctx) error {
	cid := fiber.Locals[string](c, constants.CORRELATION_ID_KEY)

	authorization := c.Get(string(constants.AUTHORIZATION_KEY))
	if authorization != "" {
		authorization = authorization[7:]
	}

	requestContext := common.BuildRequestContext(&cid, &authorization, nil, &common.UserInfo[any]{})
	c.Locals(constants.REQUEST_CONTEXT_KEY, requestContext)
	c.Locals(constants.CORRELATION_ID_KEY, requestContext.CID)

	return c.Next()
}
