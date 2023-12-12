package tokens

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	ID uuid.UUID
}

func GetStructClaims(id uuid.UUID) jwt.Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		ID: id,
	}
}

func NewSignedToken(id uuid.UUID, secretKey []byte) (string, error) {
	claim := GetStructClaims(id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, err
}

func GetMapClaims(id uuid.UUID) jwt.Claims {
	return jwt.MapClaims{
		"exp": time.Now().Add(time.Hour),
		"id":  id,
	}
}
