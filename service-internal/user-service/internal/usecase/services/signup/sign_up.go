package signup

import (
	"context"
	"strings"

	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/extractor"
	"github.com/khuongdo95/go-service/internal/adapter/dtos"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	entidentity "github.com/khuongdo95/go-service/internal/generated/ent/identity"
	entuser "github.com/khuongdo95/go-service/internal/generated/ent/user"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
	"github.com/khuongdo95/go-service/internal/usecase/services/cache"
	pass "github.com/khuongdo95/go-service/internal/usecase/services/password"
	"github.com/khuongdo95/go-service/internal/usecase/services/token"
	"github.com/khuongdo95/go-service/internal/usecase/services/transformer"
	"github.com/khuongdo95/go-service/internal/usecase/services/tx"
)

type signUp struct {
	ent       *ent.Client
	token     token.Token
	extractor extractor.Extractor
	cache     cache.UserCache
}
type SignUp interface {
	SignUp(ctx context.Context, tenantID, ip string, req *dtos.SignUpReq) (*dtos.UserRes, *response.AppError)
}

func New(ent *ent.Client, token token.Token, cache cache.UserCache) SignUp {
	return &signUp{
		ent:       ent,
		token:     token,
		extractor: extractor.New(),
		cache:     cache,
	}
}

func (s signUp) SignUp(ctx context.Context, tenantID, ip string, req *dtos.SignUpReq) (*dtos.UserRes, *response.AppError) {
	if available, errApp := s.isUsernameAvailable(ctx, req.Username); errApp != nil {
		return nil, errApp
	} else if !available {
		return nil, response.NewAppError("username already exists", constants.DuplicateData)
	}
	hashedPassword, errApp := pass.HashPassword(req.Username, req.Password)
	if errApp != nil {
		return nil, response.NewAppError("can not hash password", constants.NotFound)
	}

	var newUser *ent.User
	adminID := s.extractor.GetUserID(ctx)

	errApp = tx.WithTransaction(ctx, s.ent, func(ctx context.Context, tx tx.Tx) *response.AppError {
		var err error
		admin, err := tx.Client().User.Query().
			Where(entuser.ID(adminID)).
			Only(ctx)
		if err != nil {
			return response.ServerError(err.Error())
		}
		newUser, err = tx.Client().User.Create().
			SetTenantID(tenantID).
			SetName(req.Name).
			SetEmail(req.Email).
			SetCreatedBy(*admin.Name).
			Save(ctx)
		if err != nil {
			return response.ServerError("can not create user: " + err.Error())
		}
		_, err = tx.Client().Identity.Create().
			SetUser(newUser).
			SetUsername(strings.ToLower(req.Username)).
			SetPassword(hashedPassword).
			Save(ctx)
		if err != nil {
			return response.ServerError("can not create identity: " + err.Error())
		}

		err = tx.Client().UserIpWhiteList.Create().
			SetIPAddress(ip).
			SetUser(newUser).
			Exec(ctx)
		if err != nil {
			return response.NewAppError("can not create user ip white list", constants.InternalServerErr)
		}
		_, tokenErr := s.token.Create(tenantID, newUser.ID)
		if tokenErr != nil {
			global.Log.Error("can not create token for user", tokenErr)
			return response.NewAppError("can not create token for user", constants.InternalServerErr)
		}
		return nil
	})

	if errApp != nil {
		return nil, errApp
	}
	newUser, err := s.ent.User.Query().
		Where(entuser.ID(newUser.ID)).
		WithIdentity().
		WithIPWhiteList().
		Only(ctx)
	if err != nil {
		return nil, response.ServerError("can not get user after create: " + err.Error())
	}
	return transformer.TransformerUser(newUser), nil
}

func (s signUp) isUsernameAvailable(ctx context.Context, username string) (bool, *response.AppError) {
	existed, err := s.ent.Identity.Query().
		Where(entidentity.Username(strings.ToLower(username))).
		Exist(ctx)
	if err != nil {
		return false, response.ServerError("can not check if username existed: " + err.Error())
	}

	return !existed, nil
}
