package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type authName struct {
}
type authToken struct {
}

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		auth := request.Header.Values("Authorization")
		fmt.Println(auth)
		secretKey := os.Getenv("SECRET_KEY")
		name := auth[0]
		token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
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

		ctx := contextWithName(request.Context(), name)
		ctx = contextWithToken(request.Context(), token)
		newRequest := request.WithContext(ctx)

		handler.ServeHTTP(writer, newRequest)
	})

}

func contextWithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, authName{}, name)
}

func contextWithToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, authToken{}, token)
}

func TokenFromContext(ctx context.Context) *jwt.Token {
	name, _ := ctx.Value(authToken{}).(*jwt.Token)

	return name
}

func AccountNameFromContext(ctx context.Context) string {
	name, _ := ctx.Value(authName{}).(string)

	return name
}
