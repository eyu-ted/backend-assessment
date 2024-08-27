package domain

type LogUsecase interface {
	CreateLog(eventType string, userID string, details string, success bool) error
	GetLogs() ([]*Log, error)
}
type LogRepository interface {
	CreateLog(log *Log) error
	GetLogs() ([]*Log, error)
}
