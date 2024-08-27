// domain/log.go
package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventType string             `bson:"event_type"` // e.g., "login_attempt", "loan_submission", "loan_status_update"
	UserID    primitive.ObjectID `bson:"user_id"`
	Details   string             `bson:"details"` // Description of the event
	Success   bool               `bson:"success"` // For events like login attempts
	CreatedAt int64              `bson:"created_at"`
}
