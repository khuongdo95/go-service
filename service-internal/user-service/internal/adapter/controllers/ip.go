package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/common/utils"
	"github.com/khuongdo95/go-service/internal/adapter/dtos"
	"github.com/khuongdo95/go-service/internal/infrastructure/helper"
	"github.com/khuongdo95/go-service/internal/usecase/services"
)

type IpAccessController struct {
	ipAccessService services.IPAccessServer
}

func NewIpAccessController(ipAccessService services.IPAccessServer) *IpAccessController {
	return &IpAccessController{
		ipAccessService: ipAccessService,
	}
}

func (h *IpAccessController) CreateIP(ctx fiber.Ctx) error {
	req, err := utils.BytesToStruct[dtos.CreateIPReq](ctx.Body())
	if err != nil {
		return response.NewAppError("invalid request body", constants.ParamInvalid)
	}
	entity, errApp := h.ipAccessService.Create(ctx.Context(), req)
	if errApp != nil {
		return errApp
	}

	ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), entity))
	return nil
}

func (h *IpAccessController) GetIP(ctx fiber.Ctx) error {
	idstr := ctx.Params("id")

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return response.NewAppError("invalid id", constants.ParamInvalid)
	}
	entity, errApp := h.ipAccessService.Get(ctx.Context(), &dtos.GetIPReq{Id: id})
	if errApp != nil {
		return errApp
	}

	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), entity))
}

func (h *IpAccessController) ListIP(ctx fiber.Ctx) error {
	req := dtos.ListIPReq{}
	err := utils.MapToStruct(ctx.Queries(), &req)
	if err != nil {
		return response.NewAppError("invalid request body", constants.ParamInvalid)
	}
	entity, errApp := h.ipAccessService.List(ctx.Context(), &req)
	if errApp != nil {
		return errApp
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), entity))
}

func (h *IpAccessController) DeleteIP(ctx fiber.Ctx) error {
	idstr := ctx.Params("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return response.NewAppError("invalid id", constants.ParamInvalid)
	}
	errApp := h.ipAccessService.Delete(ctx.Context(), &dtos.DeleteIPReq{Id: id})
	if errApp != nil {
		return errApp
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), nil))
}

func (h *IpAccessController) UpdateIP(ctx fiber.Ctx) error {
	idstr := ctx.Params("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return response.NewAppError("invalid id", constants.ParamInvalid)
	}
	errApp := h.ipAccessService.Update(ctx.Context(), &dtos.UpdateIPReq{
		Id: id,
	})
	if errApp != nil {
		return errApp
	}
	return ctx.JSON(common.SuccessResponse(helper.GetHttpReqCtx(ctx), &dtos.UpdateIPRes{Success: true}))
}
