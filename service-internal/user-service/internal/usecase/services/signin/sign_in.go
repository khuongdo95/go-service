package signin

import (
	"context"
	"strings"

	"github.com/khuongdo95/go-pkg/common/constants"
	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/extractor"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	entidentity "github.com/khuongdo95/go-service/internal/generated/ent/identity"
	cache "github.com/khuongdo95/go-service/internal/usecase/services/cache"
	"github.com/khuongdo95/go-service/internal/usecase/services/limit"
	pass "github.com/khuongdo95/go-service/internal/usecase/services/password"
	"github.com/khuongdo95/go-service/internal/usecase/services/token"
)

type SignIn interface {
	SignIn(ctx context.Context, email string, password string) (*token.Tokens, *response.AppError)
}

type signIn struct {
	ent       *ent.Client
	token     token.Token
	extractor extractor.Extractor
	limit     limit.Limit
	cache     cache.UserCache
}

func New(ent *ent.Client, token token.Token, cache cache.UserCache) SignIn {
	return &signIn{
		ent:       ent,
		token:     token,
		extractor: extractor.New(),
		limit:     limit.New(ent, 5),
		cache:     cache,
	}
}

func (s signIn) SignIn(ctx context.Context, username string, password string) (*token.Tokens, *response.AppError) {
	xff := s.extractor.GetXForwardedFor(ctx)
	errApp := s.limit.CheckIP(ctx, xff)
	if errApp != nil {
		return nil, response.NewAppError("ip not in whitelist", constants.NotFound)
	}
	identityDB, err := s.ent.Identity.Query().
		Where(entidentity.Username(strings.ToLower(username))).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, response.NewAppError("username not found", constants.NotFound)
	}

	errApp = pass.CheckPassword(identityDB.Password, username, password)
	if errApp != nil {
		return nil, response.NewAppError("invalid password", constants.AcceptDenied)
	}
	tenantID := s.extractor.GetTenantID(ctx)
	tokens, errApp := s.token.Create(tenantID, identityDB.Edges.User.ID)
	if errApp != nil {
		return nil, errApp
	}
	go s.cache.Set(ctx, &cache.User{Id: identityDB.Edges.User.ID})
	return tokens, nil
}
