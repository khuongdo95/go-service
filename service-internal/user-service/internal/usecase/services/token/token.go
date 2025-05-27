package token

import (
	"github.com/khuongdo95/go-pkg/common/response"
	config "github.com/khuongdo95/go-service/internal/infrastructure/global"
	"github.com/khuongdo95/go-service/internal/usecase/services/token/signer"
)

const (
	_defaultAudience = "smm"
	_idToken         = "id_token"
	_accessToken     = "access_token"
)

type Tokens struct {
	UserID               string
	IDToken, AccessToken *signer.Token
}

type Token interface {
	AccessToken() signer.Signer
	IDToken() signer.Signer

	Create(tenantID, safeID string) (tokens *Tokens, err *response.AppError)
}

type token struct {
	idToken     signer.Signer
	accessToken signer.Signer
}

func New(accessTokenSigning *config.JwtSigning, idTokenSigning *config.JwtSigning) (Token, *response.AppError) {
	accessToken, err := signer.New(_accessToken, accessTokenSigning)
	if err != nil {
		return nil, response.ServerError(err.Error())
	}
	idToken, err := signer.New(_idToken, idTokenSigning)
	if err != nil {
		return nil, response.ServerError(err.Error())
	}

	return &token{
		idToken:     idToken,
		accessToken: accessToken,
	}, nil
}

func (t *token) AccessToken() signer.Signer {
	return t.accessToken
}

func (t *token) IDToken() signer.Signer {
	return t.idToken
}

// Create TODO audience
func (t *token) Create(tenantID, safeID string) (*Tokens, *response.AppError) {
	idToken, err := t.idToken.Create(tenantID, safeID, _defaultAudience)
	if err != nil {
		return nil, err
	}
	accessToken, err := t.accessToken.Create(tenantID, safeID, _defaultAudience)
	if err != nil {
		return nil, err
	}

	return &Tokens{safeID, idToken, accessToken}, nil
}
