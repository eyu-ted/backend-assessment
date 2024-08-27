// usecase/log_usecase.go
package usecase

import (
	"loan-tracker/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogUsecase struct {
	logRepo domain.LogRepository
}

func NewLogUsecase(logRepo domain.LogRepository) *LogUsecase {
	return &LogUsecase{
		logRepo: logRepo,
	}
}

func (u *LogUsecase) CreateLog(eventType string, userID string, details string, success bool) error {

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		// handle the error, e.g. return an error or log it
		return err
	}

	log := &domain.Log{
		EventType: eventType,
		Details:   details,
		UserID:    userIDObj,
		Success:   success,
		CreatedAt: time.Now().Unix(),
	}
	return u.logRepo.CreateLog(log)
}

func (u *LogUsecase) GetLogs() ([]*domain.Log, error) {
	return u.logRepo.GetLogs()
}
