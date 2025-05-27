package signer

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	config "github.com/khuongdo95/go-service/internal/infrastructure/global"
	"github.com/stretchr/testify/assert"
)

const (
	// You can create a new key here myid/internal/service/token/token_test.go:20
	_privateKey  = ""
	_tenantID    = "CLIENT"
	_safeID      = "SafeID"
	_issuer      = "smm"
	_type        = "access_token"
	_expiresTime = 84600 // seconds
)

func TestSign(t *testing.T) {
	//	t.Skip("Manual generate new token from private key")

	signing := &config.JwtSigning{
		PrivateKey:  _privateKey,
		ExpiresTime: _expiresTime,
		Issuer:      _issuer,
	}
	signer, err := New(_type, signing)
	assert.NoError(t, err)

	token, err := signer.Create(_tenantID, _safeID, "default")
	assert.NoError(t, err)

	fmt.Printf("JTI: %s\n", token.Id)
	fmt.Printf("TenantID: %s\n", token.TenantID)
	fmt.Printf("Token: %s\n", token.Raw)
	fmt.Printf("ExpiresAt: %v\n", token.ExpiresAt)

	tk, err := signer.Parse(token.Raw)
	assert.NoError(t, err)
	assert.True(t, tk.Valid() == nil)
	assert.Equal(t, _tenantID, token.TenantID)
	assert.Equal(t, _type, token.TokenType)
	assert.Equal(t, token.ExpiresAt, tk.ExpiresAt)
}

func TestSignSimple(t *testing.T) {
	t.Skip()

	var (
		_privateKey = ""
		_expiry     = 100 * 365 * 24 * time.Hour
		_audience   = ""
		_subject    = ""
		_issuer     = ""
	)

	privateKey, err := jwt.ParseEdPrivateKeyFromPEM([]byte(_privateKey))
	assert.NoError(t, err)

	now := time.Now()
	iat := now.UTC().Unix()
	exp := now.Add(_expiry)
	claims := jwt.StandardClaims{
		Id:        uuid.New().String(),
		ExpiresAt: exp.Unix(),
		IssuedAt:  iat,
		NotBefore: iat,
		Audience:  _audience,
		Subject:   _subject,
		Issuer:    _issuer,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(privateKey)
	assert.NoError(t, err)

	println(token)
}
