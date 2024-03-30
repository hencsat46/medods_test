package usecase

import (
	"errors"
	"log"
	"medods_test/internal/handler"
	"medods_test/internal/models"
	"medods_test/internal/pkg/jwt"
	"time"
)

type usecase struct {
	repo RepositoryInterfaces
}

type RepositoryInterfaces interface {
	InsertUser(string, string) error
	//RefreshUser(models.UserToken) error
	GetToken(string) (string, error)
}

func NewUsecase(repo RepositoryInterfaces) handler.UsecaseInterfaces {
	return &usecase{repo: repo}
}

func (u *usecase) InsertUser(userData models.UserToken) error {
	if err := u.repo.InsertUser(userData.UserId, userData.RefreshToken); err != nil {
		log.Println("Cannot insert user", err)
		return err
	}

	return nil
}

func (u *usecase) RefreshUser(userData models.UserToken) error {

	accessToken := userData.AccessToken
	refreshToken := userData.RefreshToken

	accessClaims, err := jwt.ParseAccess(accessToken)
	if err != nil {
		return err
	}

	accessId := accessClaims.UserId
	accessTime := accessClaims.Time
	accessSalt := accessClaims.Salt

	refreshClaims, err := jwt.ParseRefresh(refreshToken)
	refreshId := refreshClaims.UserId
	refreshTime := refreshClaims.Time
	refreshSalt := refreshClaims.Salt

	if time.Now().Unix() < refreshClaims.ExpTime {
		log.Println("Token is expired")
		return errors.New("token is expired")
	}

	if accessId == refreshId && accessTime == refreshTime && accessSalt == refreshSalt {
		hashedToken, err := jwt.HashRefresh(refreshToken)
		if err != nil {
			log.Println(err)
			return err
		}

		if err = u.repo.InsertUser(userData.UserId, hashedToken); err != nil {
			log.Println("Cannot insert data")
			return err
		}

	} else {
		log.Println("Tokens doesn't match")
		return errors.New("tokens doesn't match")
	}

	return nil
}
