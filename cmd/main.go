package main

import (
	"context"
	"log"
	"medods_test/internal/handler"
	"medods_test/internal/pkg/db"
	"medods_test/internal/repository"
	"medods_test/internal/usecase"
	"net/http"
	"os"
	"os/signal"

	_ "medods_test/docs"

	swagger "github.com/swaggo/http-swagger"
)

//	@title		Тестовое задание на позицию Junior Backend Developer
//	@version	1.0
//	@host		localhost:3000
//	@basepath	/

func main() {

	server := &http.Server{
		Addr: "0.0.0.0:3000",
	}

	args := os.Args[1:]

	switch len(args) {
	case 0:
		os.Setenv("SECRET", "secretkey")
		os.Setenv("MONGODB_URL", "medods_db")
	case 1:
		os.Setenv("MONGODB_URL", args[0])
	case 2:
		os.Setenv("MONGODB_URL", args[0])
		os.Setenv("SECRET", args[1])
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conn := db.InitMongo()
	repo := repository.NewRepository(conn)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)
	http.HandleFunc("POST /create", handler.CreateTokens)
	http.HandleFunc("POST /refresh", handler.Refresh)
	http.HandleFunc("GET /swagger/*", swagger.WrapHandler)

	go func(ctx context.Context) {
		<-ctx.Done()
		server.Close()
	}(ctx)
	log.Println("Server started on port :3000...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
