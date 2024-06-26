package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	customErrors "medods_test/internal/errors"
	"medods_test/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		tokenTime,
		salt,
		golangJwt.RegisteredClaims{
			ExpiresAt: golangJwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
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

func ParseAccess(inputToken string) (jwtAccessClaims, error) {
	token, err := golangJwt.ParseWithClaims(inputToken, &jwtAccessClaims{}, func(t *golangJwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {

			return jwtAccessClaims{}, customErrors.ErrAccessToken
		}
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return jwtAccessClaims{}, err
		}
	}

	if claims, ok := token.Claims.(*jwtAccessClaims); ok {
		return *claims, nil
	} else {
		return jwtAccessClaims{}, errors.New("cannot parse access token")
	}

}

func ParseRefresh(inputToken string) (refreshClaims, error) {

	decodedToken, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		log.Println("Cannot decode token", err)
		return refreshClaims{}, customErrors.ErrRefreshToken
	}

	unmarshaledToken := refreshClaims{}

	if err := json.Unmarshal(decodedToken, &unmarshaledToken); err != nil {
		log.Println("Cannot unmarshal token", err)
		return refreshClaims{}, err
	}

	return unmarshaledToken, nil

}

func createRefreshToken(customTime, expTime int64, userId, salt, secret string) (string, error) {
	refreshToken := refreshClaims{UserId: userId, Salt: salt, Secret: secret, Time: customTime, ExpTime: expTime}

	tokenString, err := json.Marshal(refreshToken)
	if err != nil {
		log.Println("Cannot create refresh token: ", err)
		return "", err
	}

	encodedToken := base64.StdEncoding.EncodeToString(tokenString)

	return encodedToken, nil
}

func HashRefresh(token string) (string, error) {
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(token)); err != nil {
		log.Println("Cannot hash token", err)
		return "", err
	}
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
