package limit

import (
	"context"
	"strings"

	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	useripwhitelist "github.com/khuongdo95/go-service/internal/generated/ent/useripwhitelist"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
	"go.uber.org/zap"
)

type Limit interface {
	CheckIP(ctx context.Context, xff string) *response.AppError
}

type dayLimit struct {
	ent *ent.Client
	max uint32
}

func New(ent *ent.Client, max uint32) Limit {
	return &dayLimit{
		ent: ent,
		max: max,
	}
}

func (l *dayLimit) CheckIP(ctx context.Context, xff string) *response.AppError {
	if len(xff) == 0 {
		return nil
	}
	ip := strings.TrimSpace(strings.Split(xff, ",")[0])
	existIp, err := l.ent.UserIpWhiteList.Query().Where(useripwhitelist.IPAddress(ip)).Exist(ctx)
	if err != nil {
		return response.AccessDenied("can not check ip in whitelist")
	}

	if !existIp {
		global.Log.Error("ip not in whitelist", err, zap.String("ip", ip))
		return response.NewAppError("ip not in whitelist", constants.AcceptDenied)
	}

	return nil
}
