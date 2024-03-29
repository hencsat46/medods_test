package jwt

import (
	"errors"
	"fmt"
	"log"
	"medods_test/internal/models"
	"net/http"
	"os"
	"time"

	golangJwt "github.com/golang-jwt/jwt/v5"
)

type jwtAccessClaims struct {
	UserId string
	golangJwt.RegisteredClaims
}

func CreateAccessToken(userId string) (string, error) {

	claims := jwtAccessClaims{
		userId,
		golangJwt.RegisteredClaims{
			ExpiresAt: golangJwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := golangJwt.NewWithClaims(golangJwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("Cannot sign token", err)
		return "", err
	}

	return tokenString, nil
}

func ValidationAccessJWT(innerFunc func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(write http.ResponseWriter, read *http.Request) {
		if read.Header["Token"] != nil {

			token, err := golangJwt.ParseWithClaims(read.Header["Token"][0], &jwtAccessClaims{}, func(token *golangJwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				if errors.Is(err, golangJwt.ErrTokenExpired) {
					fmt.Fprint(write, models.Response{Status: 401, Payload: "Token is Expired"})
					return
				}
				log.Println("Token Error", err)
			}

			if token.Valid {
				innerFunc(write, read)
			}
		}
	})
}
