package repository

import (
	"context"
	"medods_test/internal/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	mongoConnection *mongo.Client
}

type mongoData struct {
	UserId      string
	HashedToken string
}

func NewRepository(conn *mongo.Client) usecase.RepositoryInterfaces {
	return &repository{mongoConnection: conn}
}

func (r *repository) InsertUser(userId, token string) error {

	insertData := mongoData{UserId: userId, HashedToken: token}

	if _, err := r.mongoConnection.Database("medods_db").Collection("tokens").InsertOne(context.Background(), insertData); err != nil {
		return err
	}

	return nil
}

// func (r *repository) RefreshUser(data models.UserToken) error {

// }

func (r *repository) GetToken(userId string) (string, error) {
	filter := bson.D{{"UserId", userId}}
	result := mongoData{}

	if err := r.mongoConnection.Database("medods_db").Collection("tokens").FindOne(context.TODO(), filter).Decode(&result); err != nil {
		return "", err
	}

	return result.HashedToken, nil

}
