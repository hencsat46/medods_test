package handler

import (
	"encoding/json"
	"fmt"
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
		response := models.Response{Status: 401, Payload: "JSON error"}
		Response(w, response)
		return
	}

	accessToken, refreshToken, err := jwt.CreateTokens(request.UserId)
	if err != nil {
		response := models.Response{Status: 500, Payload: "Create token error"}
		Response(w, response)
		return
	}

	data := models.UserToken{UserId: request.UserId, AccessToken: accessToken, RefreshToken: refreshToken}

	if err = h.usecase.InsertUser(data); err != nil {
		response := models.Response{Status: 501, Payload: "Internal Server Error"}
		Response(w, response)
		return
	}

	response := models.Response{Status: 200, Payload: fmt.Sprintf("Access Token: %v\nRefreshToken: %v", accessToken, refreshToken)}
	Response(w, response)

}

func (h *handler) Something(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {

}

func Response(writer http.ResponseWriter, responseData models.Response) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(responseData.Status)
	fmt.Fprint(writer, responseData)
}
