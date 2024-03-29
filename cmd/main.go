package main

import (
	"medods_test/internal/handler"
	"medods_test/internal/pkg/db"
	"medods_test/internal/pkg/env"
	"medods_test/internal/pkg/jwt"
	"net/http"
)

func main() {
	env.Init()
	db.InitMongo()

	handler := handler.NewHandler(nil)
	http.HandleFunc("POST /create", handler.CreateTokens)
	http.HandleFunc("GET /validate", jwt.ValidationAccessJWT(handler.Something))

	http.ListenAndServe(":3000", nil)
}
