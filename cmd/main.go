package main

import (
	"context"
	"medods_test/internal/handler"
	"medods_test/internal/pkg/db"
	"medods_test/internal/pkg/env"
	"medods_test/internal/repository"
	"medods_test/internal/usecase"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	server := &http.Server{
		Addr: ":3000",
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	env.Init()
	conn := db.InitMongo()
	repo := repository.NewRepository(conn)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)
	http.HandleFunc("POST /create", handler.CreateTokens)
	http.HandleFunc("POST /refresh", handler.Refresh)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				server.Close()

				return
			}
		}
	}(ctx)

	server.ListenAndServe()

}
