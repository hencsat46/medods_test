package repository

import (
	"context"
	"errors"
	"log"
	"medods_test/internal/models"
	"medods_test/internal/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *repository) InsertUser(user models.UserToken) error {

	insertData := mongoData{UserId: user.UserId, HashedToken: user.RefreshToken}

	if _, err := r.mongoConnection.Database("medods_db").Collection("tokens").InsertOne(context.Background(), insertData); err != nil {
		return err
	}

	return nil
}

func (r *repository) CheckUser(userId string) (bool, error) {
	filter := bson.D{{"userid", userId}}
	opts := options.FindOne().SetProjection(bson.D{{"hashedtoken", false}})
	result := mongoData{}

	if err := r.mongoConnection.Database("medods_db").Collection("tokens").FindOne(context.TODO(), filter, opts).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			return false, nil
		}
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (r *repository) UpdateUser(data models.UserToken) error {
	filter := bson.D{{"userid", data.UserId}}
	update := bson.D{{"$set", bson.D{{"hashedtoken", data.RefreshToken}}}}

	_, err := r.mongoConnection.Database("medods_db").Collection("tokens").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetToken(userId string) (string, error) {
	filter := bson.D{{"userid", userId}}
	result := mongoData{}

	if err := r.mongoConnection.Database("medods_db").Collection("tokens").FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Println(err)
		return "", err
	}

	return result.HashedToken, nil

}
