package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/common/utils"
	"github.com/khuongdo95/go-service/internal/adapter/dtos"
	"github.com/khuongdo95/go-service/internal/infrastructure/helper"
	"github.com/khuongdo95/go-service/internal/usecase/services"
)

type (
	UserController struct {
		userService services.UserServiceIf
	}
)

func NewUserController(userService services.UserServiceIf) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) Me(ctx fiber.Ctx) error {
	userID := helper.GetHttpUserCtx(ctx).ID

	dataUser, err := h.userService.Get(ctx.Context(), userID)
	if err != nil {
		return err
	}
	ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), dataUser))
	return nil
}

func (h *UserController) Get(ctx fiber.Ctx) error {
	dataUser, err := h.userService.Get(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}
	ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), &dtos.UserRes{
		Id:   dataUser.ID,
		Name: *dataUser.Name,
	}))
	return nil
}

func (h *UserController) SignUp(ctx fiber.Ctx) error {
	dto, err := utils.BytesToStruct[dtos.SignUpReq](ctx.Body())
	if err != nil {
		return response.NewAppError("invalid request body", constants.ParamInvalid)
	}

	tokenUser, errApp := h.userService.SignUp(ctx.Context(), dto.Username, dto.Password)
	if errApp != nil {
		return errApp
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), &dtos.SignUpRes{
		Token: &dtos.TokenInfo{
			AccessToken: tokenUser.AccessToken.Raw,
			IdToken:     tokenUser.IDToken.Raw,
			UserId:      tokenUser.UserID,
		},
	}))
}

func (h *UserController) SignIn(ctx fiber.Ctx) error {
	req, err := utils.BytesToStruct[dtos.SignInReq](ctx.Body())
	if err != nil {
		return response.NewAppError("invalid request body", constants.ParamInvalid)
	}

	tokenUser, errApp := h.userService.SignIn(ctx.Context(), req.Username, req.Password)
	if errApp != nil {
		return errApp
	}

	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), &dtos.SignInRes{
		Token: &dtos.TokenInfo{
			AccessToken: tokenUser.AccessToken.Raw,
			IdToken:     tokenUser.IDToken.Raw,
			UserId:      tokenUser.UserID,
		},
	}))
}

func (h *UserController) Update(ctx fiber.Ctx) error {
	req, err := utils.BytesToStruct[dtos.UpdateUserReq](ctx.Body())
	if err != nil {
		return response.NewAppError("invalid request body", constants.ParamInvalid)
	}
	dataUser, errApp := h.userService.Update(ctx.Context(), req)
	if errApp != nil {
		return errApp
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), dataUser))
}

func (h *UserController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.NewAppError("user id is required", constants.ParamInvalid)
	}
	err := h.userService.Delete(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), nil))
}

func (h *UserController) Create(ctx fiber.Ctx) error {
	req, err := utils.BytesToStruct[dtos.CreateUserReq](ctx.Body())
	if err != nil {
		return response.NewAppError("invalid request body", constants.ParamInvalid)
	}
	dataUser, errApp := h.userService.Create(ctx.Context(), req)
	if errApp != nil {
		return errApp
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), dataUser))
}

func (h *UserController) List(ctx fiber.Ctx) error {
	req := dtos.ListUserReq{}
	err := utils.MapToStruct(ctx.Queries(), &req)
	if err != nil {
		return response.NewAppError("invalid request query", constants.ParamInvalid)
	}

	dataUser, errApp := h.userService.List(ctx.Context(), &req)
	if errApp != nil {
		return errApp
	}

	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), dataUser))
}
