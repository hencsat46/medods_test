package jwt

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"medods_test/internal/models"
	"net/http"
	"os"
	"time"

	golangJwt "github.com/golang-jwt/jwt/v5"
)

type jwtAccessClaims struct {
	UserId string
	Time   int64
	Salt   string
	golangJwt.RegisteredClaims
}

type refreshClaims struct {
	UserId  string
	Salt    string
	Secret  string
	Time    int64
	ExpTime int64
	IsUsed  bool
}

func CreateTokens(userId string) (string, string, error) {

	tokenTime := time.Now().Unix()
	salt, err := createSalt()
	if err != nil {
		log.Println("Cannot create salt", err)
	}

	refreshToken, err := createRefreshToken(tokenTime, time.Now().Add(time.Minute*5).Unix(), userId, salt, os.Getenv("SECRET"))
	if err != nil {
		return "", "", err
	}

	claims := jwtAccessClaims{
		userId,
		time.Now().Unix(),
		salt,
		golangJwt.RegisteredClaims{
			ExpiresAt: golangJwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := golangJwt.NewWithClaims(golangJwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("Cannot sign token", err)
		return "", "", err
	}

	return tokenString, refreshToken, nil
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

func createRefreshToken(customTime, expTime int64, userId, salt, secret string) (string, error) {
	refreshToken := refreshClaims{UserId: userId, Salt: salt, Secret: secret, Time: customTime}

	tokenString, err := json.Marshal(refreshToken)
	if err != nil {
		log.Println("Cannot create refresh token: ", err)
		return "", err
	}

	hasher := sha512.New()
	hasher.Write(tokenString)
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}

func createSalt() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
