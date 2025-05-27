package transformer

import (
	"github.com/khuongdo95/go-service/internal/adapter/dtos"
	"github.com/khuongdo95/go-service/internal/generated/ent"
)

func TransformerUser(user *ent.User) *dtos.UserRes {
	var ip string
	if len(user.Edges.IPWhiteList) > 0 {
		ip = user.Edges.IPWhiteList[0].IPAddress
	}

	return &dtos.UserRes{
		Id:        user.ID,
		Username:  user.Edges.Identity[0].Username,
		Name:      *user.Name,
		Email:     *user.Email,
		IpAddress: ip,
		Metadata: &dtos.ModifiedEntity{
			CreatedAt: &user.CreatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.CreatedBy,
		},
	}
}

func TransformerIP(ip *ent.IPAccessControl) *dtos.IPAccess {
	return &dtos.IPAccess{
		Id:            ip.ID,
		IpAddress:     ip.IPAddress,
		RabbitmqLive:  ip.RabbitmqLive,
		RabbitmqStage: ip.RabbitmqStage,
		TcpLive:       ip.TCPLive,
		TcpStage:      ip.TCPStage,
		Active:        ip.Active,
		Metadata: &dtos.ModifiedEntity{
			CreatedAt: &ip.CreatedAt,
			UpdatedAt: ip.UpdatedAt,
		},
	}
}
