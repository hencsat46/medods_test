package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	customErrors "medods_test/internal/errors"
	"medods_test/internal/models"
	"medods_test/internal/pkg/jwt"
	"net/http"
)

type handler struct {
	usecase UsecaseInterfaces
}

type UsecaseInterfaces interface {
	InsertUser(models.UserToken) error
	RefreshUser(models.UserToken) (models.UserToken, error)
}

func NewHandler(usecase UsecaseInterfaces) *handler {
	return &handler{usecase: usecase}
}

func (h *handler) CreateTokens(w http.ResponseWriter, r *http.Request) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || len(request.UserId) == 0 {
		response := models.Response{Status: 400, Payload: "JSON error"}
		h.Response(w, response)
		return
	}

	accessToken, refreshToken, err := jwt.CreateTokens(request.UserId)
	if err != nil {
		response := models.Response{Status: 500, Payload: "Create token error"}
		h.Response(w, response)
		return
	}

	data := models.UserToken{UserId: request.UserId, AccessToken: accessToken, RefreshToken: refreshToken}

	if err = h.usecase.InsertUser(data); err != nil {
		if errors.Is(err, customErrors.ErrUserExists) {
			response := models.Response{Status: 400, Payload: "This user already exists"}
			h.Response(w, response)
			return
		}
		response := models.Response{Status: 500, Payload: "Internal Server Error"}
		h.Response(w, response)
		return
	}

	response := models.Response{Status: 200, Payload: struct {
		AccessToken  string
		RefreshToken string
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}}
	h.Response(w, response)

}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || len(request.UserId) == 0 {
		response := models.Response{Status: 400, Payload: "JSON error"}
		h.Response(w, response)
		return
	}

	userData := models.UserToken{UserId: request.UserId, AccessToken: r.Header["Access"][0], RefreshToken: r.Header["Refresh"][0]}

	data, err := h.usecase.RefreshUser(userData)

	if err != nil {

		if errors.Is(err, customErrors.ErrTokenExpired) {
			response := models.Response{Status: 401, Payload: "Refresh token is expired"}
			h.Response(w, response)
			return
		}

		if errors.Is(err, customErrors.ErrTokensMatch) {
			response := models.Response{Status: 401, Payload: "Tokens doesn't match"}
			h.Response(w, response)
			return
		}

		if errors.Is(err, customErrors.ErrRefreshToken) {
			response := models.Response{Status: 400, Payload: "Invalid refresh token"}
			h.Response(w, response)
			return
		}

		if errors.Is(err, customErrors.ErrAccessToken) {
			response := models.Response{Status: 400, Payload: "Invalid access token"}
			h.Response(w, response)
			return
		}

		response := models.Response{Status: 500, Payload: "Internal server error"}
		h.Response(w, response)
		return
	}

	response := models.Response{Status: 200, Payload: struct {
		AccessToken  string
		RefreshToken string
	}{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}}

	h.Response(w, response)

}

func (h *handler) Response(writer http.ResponseWriter, responseData models.Response) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(responseData.Status)

	jsonString, err := json.Marshal(responseData)
	if err != nil {
		log.Println("Cannot marshal json")
		fmt.Fprint(writer, responseData)
	}

	fmt.Fprint(writer, string(jsonString))

}
