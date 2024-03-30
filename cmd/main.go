package main

import (
	"medods_test/internal/handler"
	"medods_test/internal/pkg/db"
	"medods_test/internal/pkg/env"
	"medods_test/internal/repository"
	"medods_test/internal/usecase"
	"net/http"
)

func main() {
	env.Init()
	conn := db.InitMongo()
	repo := repository.NewRepository(conn)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)
	http.HandleFunc("POST /create", handler.CreateTokens)
	http.HandleFunc("GET /test", handler.Something)

	http.ListenAndServe(":3000", nil)
}
