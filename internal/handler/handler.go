package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"medods_test/internal/models"
	"medods_test/internal/pkg/jwt"
	"net/http"
)

type handler struct {
	usecase UsecaseInterfaces
}

type UsecaseInterfaces interface {
	InsertUser(models.UserToken) error
	RefreshUser(models.UserToken) error
}

func NewHandler(usecase UsecaseInterfaces) *handler {
	return &handler{usecase: usecase}
}

func (h *handler) CreateTokens(w http.ResponseWriter, r *http.Request) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || len(request.UserId) == 0 {
		fmt.Fprint(w, models.Response{Status: 401, Payload: "JSON error"})
		return
	}

	log.Println(request.UserId)

	accessToken, refreshToken, err := jwt.CreateTokens(request.UserId)
	if err != nil {
		fmt.Fprint(w, models.Response{Status: 500, Payload: "Create token error"})
		return
	}

	data := models.UserToken{UserId: request.UserId, AccessToken: accessToken, RefreshToken: refreshToken}

	if err = h.usecase.InsertUser(data); err != nil {
		fmt.Fprint(w, models.Response{Status: 501, Payload: "Internal Server Error"})
	}

	fmt.Fprint(w, models.Response{Status: 200, Payload: fmt.Sprintf("Access Token: %v\nRefreshToken: %v", accessToken, refreshToken)})

}

func (h *handler) Something(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, models.Response{Status: 200, Payload: "Token Ok"})
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {

}
