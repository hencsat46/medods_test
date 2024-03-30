package handler

import (
	"encoding/json"
	"fmt"
	"medods_test/internal/models"
	"medods_test/internal/pkg/jwt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	usecase UsecaseInterfaces
}

type UsecaseInterfaces interface {
	CreateTokens(ctx echo.Context) error
	AuthToken(ctx echo.Context) error
}

func NewHandler(usecase UsecaseInterfaces) *handler {
	return &handler{usecase: usecase}
}

func (h *handler) CreateTokens(w http.ResponseWriter, r *http.Request) {
	requestBody := make(map[string]string, 1)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprint(w, models.Response{Status: 401, Payload: "JSON error"})
		return
	}

	accessToken, refreshToken, err := jwt.CreateTokens(requestBody["UserId"])
	if err != nil {
		fmt.Fprint(w, models.Response{Status: 500, Payload: "Create token error"})
		return
	}

	fmt.Fprint(w, models.Response{Status: 200, Payload: fmt.Sprintf("Access Token: %v\nRefreshToken: %v", accessToken, refreshToken)})

}

func (h *handler) Something(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, models.Response{Status: 200, Payload: "Token Ok"})
}
