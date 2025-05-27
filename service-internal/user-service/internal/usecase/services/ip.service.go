package services

import (
	"context"

	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/extractor"
	"github.com/khuongdo95/go-service/internal/adapter/dtos"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	"github.com/khuongdo95/go-service/internal/generated/ent/ipaccesscontrol"
	"github.com/khuongdo95/go-service/internal/usecase/services/transformer"
)

type IpAccessServer struct {
	ent *ent.Client
	ext extractor.Extractor
}

type IPAccessServer interface {
	IsAllowed(ctx context.Context, ipAddress string) bool
	Create(ctx context.Context, req *dtos.CreateIPReq) (*dtos.IPAccess, *response.AppError)
	Update(ctx context.Context, req *dtos.UpdateIPReq) *response.AppError
	Delete(ctx context.Context, req *dtos.DeleteIPReq) *response.AppError
	List(ctx context.Context, req *dtos.ListIPReq) (*dtos.ListIPRes, *response.AppError)
	Get(ctx context.Context, req *dtos.GetIPReq) (*dtos.IPAccess, *response.AppError)
}

func NewIpAccess(ent *ent.Client) *IpAccessServer {
	return &IpAccessServer{
		ent: ent,
		ext: extractor.New(),
	}
}

func (i *IpAccessServer) IsAllowed(ctx context.Context, ipAddress string) bool {
	exists, err := i.ent.IPAccessControl.Query().
		Where(
			ipaccesscontrol.IPAddress(ipAddress),
			ipaccesscontrol.Active(true),
		).
		Exist(ctx)
	if err != nil {
		return false
	}
	return exists
}

func (i *IpAccessServer) Create(ctx context.Context, Req *dtos.CreateIPReq) (*dtos.IPAccess, *response.AppError) {
	created, err := i.ent.IPAccessControl.Create().
		SetIPAddress(Req.IpAddress).
		SetRabbitmqLive(Req.RabbitmqLive).
		SetRabbitmqStage(Req.RabbitmqStage).
		SetTCPLive(Req.TcpLive).
		SetTCPStage(Req.TcpStage).
		SetCreatedBy("System").
		Save(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}
	return transformer.TransformerIP(created), nil
}

func (i *IpAccessServer) Update(ctx context.Context, req *dtos.UpdateIPReq) *response.AppError {
	err := i.ent.IPAccessControl.UpdateOneID(req.Id).
		SetNillableIPAddress(req.IpAddress).
		SetNillableRabbitmqLive(req.RabbitmqLive).
		SetNillableRabbitmqStage(req.RabbitmqStage).
		SetNillableTCPLive(req.TcpLive).
		SetNillableTCPStage(req.TcpStage).
		Exec(ctx)
	if err != nil {
		return response.ConvertDatabaseError(err)
	}
	return nil
}

func (i *IpAccessServer) Delete(ctx context.Context, req *dtos.DeleteIPReq) *response.AppError {
	if err := i.ent.IPAccessControl.DeleteOneID(req.Id).Exec(ctx); err != nil {
		return response.ConvertDatabaseError(err)
	}
	return nil
}

func (i *IpAccessServer) List(ctx context.Context, req *dtos.ListIPReq) (*dtos.ListIPRes, *response.AppError) {
	var IpAccessList []*dtos.IPAccess

	query := i.ent.IPAccessControl.Query()
	total, err := query.Count(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}

	entities, err := query.
		Limit(int(req.Pagination.PageSize)).
		Offset(int(req.Pagination.PageIndex * req.Pagination.PageSize)).
		All(ctx)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}

	for _, entity := range entities {
		IpAccessList = append(IpAccessList, transformer.TransformerIP(entity))
	}

	return &dtos.ListIPRes{
		Data: IpAccessList,
		Pagination: &dtos.PaginationRes{
			PageSize:  req.Pagination.PageSize,
			PageIndex: req.Pagination.PageIndex,
			Total:     int32(total),
		},
	}, nil
}

func (i *IpAccessServer) Get(ctx context.Context, req *dtos.GetIPReq) (*dtos.IPAccess, *response.AppError) {
	ipControl, err := i.ent.IPAccessControl.Get(ctx, req.Id)
	if err != nil {
		return nil, response.ConvertDatabaseError(err)
	}

	return transformer.TransformerIP(ipControl), nil
}
