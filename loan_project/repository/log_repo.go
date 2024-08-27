// repository/log_repository_mongo.go
package repository

import (
	"context"
	"loan-tracker/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoLogRepository struct {
	db *mongo.Collection
}

func NewLogRepository(db *mongo.Database) domain.LogRepository {
	return &mongoLogRepository{
		db: db.Collection("logs"),
	}
}

func (r *mongoLogRepository) CreateLog(log *domain.Log) error {
	_, err := r.db.InsertOne(context.TODO(), log)
	return err
}

func (r *mongoLogRepository) GetLogs() ([]*domain.Log, error) {
	var logs []*domain.Log
	cursor, err := r.db.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var log domain.Log
		if err = cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, cursor.Err()
}
