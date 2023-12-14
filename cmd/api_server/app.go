package main

import (
	"bank-api/internal/di"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	container := di.NewContainer(ctx)
	defer container.Close()
	err := http.ListenAndServe(":8080", container.HTTPRouter())
	if err != nil {
		log.Fatal(err)
	}
}
