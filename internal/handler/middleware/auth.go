package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type authToken struct {
}

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")
		secretKey := os.Getenv("SECRET_KEY")
		token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Недопустимый метод подписи: %v", token.Header["alg"])
			}

			// Возвращаем секретный ключ для проверки подписи
			return []byte(secretKey), nil
		})

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
