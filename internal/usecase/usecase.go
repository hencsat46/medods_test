package usecase

import (
	"log"
	customErrors "medods_test/internal/errors"
	"medods_test/internal/handler"
	"medods_test/internal/models"
	"medods_test/internal/pkg/jwt"
	"time"
)

type usecase struct {
	repo RepositoryInterfaces
}

type RepositoryInterfaces interface {
	InsertUser(models.UserToken) error
	UpdateUser(models.UserToken) error
	GetToken(string) (string, error)
	CheckUser(string) (bool, error)
}

func NewUsecase(repo RepositoryInterfaces) handler.UsecaseInterfaces {
	return &usecase{repo: repo}
}

func (u *usecase) InsertUser(userData models.UserToken) error {

	isExists, err := u.repo.CheckUser(userData.UserId)

	if err != nil {
		log.Println("Cannot get data", err)
		return err
	}

	if isExists {
		hashedToken, err := jwt.HashRefresh(userData.RefreshToken)
		if err != nil {
			log.Println("Cannot hash refresh token")
			return err
		}
		userData.RefreshToken = hashedToken
		if err = u.repo.UpdateUser(userData); err != nil {
			log.Println("Cannot insert user", err)
			return err
		}
		return nil
	}

	hashedToken, err := jwt.HashRefresh(userData.RefreshToken)
	if err != nil {
		log.Println("Cannot hash refresh token", err)
		return err
	}

	userData.RefreshToken = hashedToken

	if err := u.repo.InsertUser(userData); err != nil {
		log.Println("Cannot insert user", err)
		return err
	}

	return nil
}

func (u *usecase) RefreshUser(userData models.UserToken) (models.UserToken, error) {

	accessToken := userData.AccessToken
	refreshToken := userData.RefreshToken

	accessClaims, err := jwt.ParseAccess(accessToken)
	if err != nil {
		log.Println(err)
		return models.UserToken{}, err
	}

	accessId := accessClaims.UserId
	accessTime := accessClaims.Time
	accessSalt := accessClaims.Salt

	refreshClaims, err := jwt.ParseRefresh(refreshToken)
	if err != nil {
		log.Println("Cannot parse token", err)

		return models.UserToken{}, err
	}

	refreshId := refreshClaims.UserId
	refreshTime := refreshClaims.Time
	refreshSalt := refreshClaims.Salt

	if refreshClaims.ExpTime < time.Now().Unix() {
		log.Println("Refresh token is expired")
		return models.UserToken{}, customErrors.ErrTokenExpired
	}

	if accessId == refreshId && accessTime == refreshTime && accessSalt == refreshSalt {

		hashedToken, err := jwt.HashRefresh(refreshToken)
		if err != nil {
			log.Println("Cannot hash refresh token", err)
			return models.UserToken{}, err
		}

		oldToken, err := u.repo.GetToken(userData.UserId)
		if err != nil {
			log.Println("Cannot get old token", err)
		}

		if oldToken != hashedToken {
			log.Println("Refresh tokens doesn't match")
			return models.UserToken{}, customErrors.ErrTokensMatch
		}

		newAccessToken, newRefreshToken, err := jwt.CreateTokens(userData.UserId)
		if err != nil {
			log.Println("Cannot create tokens", err)
			return models.UserToken{}, err
		}

		hashedToken, err = jwt.HashRefresh(newRefreshToken)
		if err != nil {
			log.Println("Cannot hash refresh token", err)
			return models.UserToken{}, err
		}

		newData := models.UserToken{UserId: userData.UserId, AccessToken: newAccessToken, RefreshToken: hashedToken}

		if err = u.repo.UpdateUser(newData); err != nil {
			log.Println("Cannot insert data")
			return models.UserToken{}, err
		}

		newData.RefreshToken = newRefreshToken
		return newData, nil

	} else {
		log.Println("Tokens doesn't match")
		return models.UserToken{}, customErrors.ErrTokensMatch
	}

}
