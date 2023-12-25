package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authToken struct {
}

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")
		secretKey := os.Getenv("SECRET_KEY")
		token, err := ValidateAndParseToken(auth, secretKey)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := contextWithToken(request.Context(), token)
		newRequest := request.WithContext(ctx)

		handler.ServeHTTP(writer, newRequest)
	})
}

func contextWithToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, authToken{}, token)
}

func TokenFromContext(ctx context.Context) *jwt.Token {
	name, _ := ctx.Value(authToken{}).(*jwt.Token)

	return name
}

func ValidateAndParseToken(auth string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Недопустимый метод подписи: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	return token, err
}

func validateTokenExpiration(token *jwt.Token) error {
	if !token.Valid {
		return fmt.Errorf("Токен не действителен")
	}

	expirationTime, ok := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return fmt.Errorf("Время истечения токена не найдено")
	}

	expiration := time.Unix(int64(expirationTime), 0)

	if time.Now().After(expiration) {
		return fmt.Errorf("Токен истек")
	}

	return nil
}
