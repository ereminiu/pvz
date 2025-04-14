package role_token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	signingKey = `cGcfFMIo+19@#-_12()~\(*V*)/~cEC]\\BaRjohi1940194WQnWHCFRkrqpn[jZLGdFCYreK]k\192018340183`
	tokenTTL   = 96 * time.Hour
	adminRole  = "admin"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserRole string `json:"user_role"`
}

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		adminRole,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseRole(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return "", nil
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserRole, nil
}

func CheckRole(role string) bool {
	return role == adminRole
}
