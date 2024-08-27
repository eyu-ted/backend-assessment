// repository/user_repository_mongo.go
package repository

import (
	"context"
	"errors"
	"loan-tracker/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &mongoUserRepository{
		db: db.Collection("users"),
	}
}

func (r *mongoUserRepository) CreateUser(user *domain.User) error {
	_, err := r.db.InsertOne(context.TODO(), user)
	return err
}

func (r *mongoUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"email": email}
	err := r.db.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) UpdateUser(user *domain.User) error {
	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": user}
	_, err := r.db.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *mongoUserRepository) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{
		"_id":               userID,
		"is_email_verified": true,
	}
	err := r.db.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (r *mongoUserRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	cursor, err := r.db.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user domain.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *mongoUserRepository) DeleteUser(userID string) error {
	filter := bson.M{"_id": userID}
	result, err := r.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
