package signer

import (
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/khuongdo95/go-pkg/common/response"
	config "github.com/khuongdo95/go-service/internal/infrastructure/global"
)

type Claims struct {
	TenantID  string `json:"tenant_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

type Token struct {
	*Claims
	Raw string
}

type Signer interface {
	Create(tenantID, safeID string, audience string) (token *Token, err *response.AppError)
	Parse(token string) (*Token, *response.AppError)
}

type signer struct {
	tokenType  string
	privateKey ed25519.PrivateKey
	expiry     time.Duration
	issuer     string
}

func New(tokenType string, jwtSigning *config.JwtSigning) (Signer, *response.AppError) {
	privateKey, err := jwt.ParseEdPrivateKeyFromPEM([]byte(jwtSigning.PrivateKey))
	if err != nil {
		return nil, response.ServerError("failed to parse private key: " + err.Error())
	}
	expiry := time.Duration(jwtSigning.ExpiresTime)

	return &signer{
		tokenType:  tokenType,
		privateKey: privateKey.(ed25519.PrivateKey),
		expiry:     expiry,
		issuer:     jwtSigning.Issuer,
	}, nil
}

func (t *signer) Create(tenantID, safeID string, audience string) (*Token, *response.AppError) {
	now := time.Now()
	iat := now.UTC().Unix()
	exp := now.Add(t.expiry * time.Second)
	id := uuid.New().String()
	claims := &Claims{
		TokenType: t.tokenType,
		TenantID:  tenantID,
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			ExpiresAt: exp.Unix(),
			IssuedAt:  iat,
			NotBefore: iat,
			Audience:  audience,
			Subject:   safeID,
			Issuer:    t.issuer,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(t.privateKey)

	return &Token{
		Raw:    token,
		Claims: claims,
	}, response.ServerError(err.Error())
}

func (t *signer) Parse(token string) (*Token, *response.AppError) {
	tk, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return t.privateKey.Public(), nil
	})
	if err != nil {
		return nil, response.Unauthorized("failed to parse token: " + err.Error())
	}

	if claims, ok := tk.Claims.(*Claims); ok && tk.Valid {
		return &Token{
			Claims: claims,
			Raw:    token,
		}, nil
	} else {
		return nil, response.Unauthorized("invalid token claims or token is not valid")
	}
}
