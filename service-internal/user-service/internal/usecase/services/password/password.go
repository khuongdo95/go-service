package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/khuongdo95/go-pkg/common/response"
	"golang.org/x/crypto/bcrypt"
)

const (
	_defaultPasswordLength = 10
	_passwordChars         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var _cost = bcrypt.DefaultCost

func Cost(cost int32) {
	// _cost value < bcrypt.MinCost or > bcrypt.MaxCost will be override by bcrypt while executing GenerateFromPassword
	_cost = int(cost)
}

func CreateMyID(username, password string) string {
	return fmt.Sprintf("%s:%s", username, password)
}

func HashPassword(username, password string) (string, *response.AppError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(CreateMyID(username, password)), _cost)
	if err != nil {
		return "", response.ServerError(err.Error())
	}

	return string(bytes), nil
}

func CheckPassword(hashedPassword, username, password string) *response.AppError {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(CreateMyID(username, password)))
	if err != nil {
		return response.AccessDenied("invalid password")
	}
	return nil
}

func GeneratePassword(username string) (raw, hashed string, err *response.AppError) {
	var b strings.Builder
	chars := []rune(_passwordChars)

	for i := 0; i < _defaultPasswordLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", "", response.ServerError(err.Error())
		}
		b.WriteRune(chars[n.Int64()])
	}

	raw = b.String()
	hashed, err = HashPassword(username, raw)
	return
}
