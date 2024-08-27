package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	GetUserByID(userID string) (*User, error)
	GetAllUsers() ([]*User, error)
	DeleteUser(userID string) error
}
type UserUsecase interface {
	RegisterUser(user *User) error
	VerifyEmail(token, email string) error
	LoginUser(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	GetUserProfile(userID string) (*User, error)
	RequestPasswordReset(email string) error
	GetAllUsers() ([]*User, error)
	DeleteUserByID(userID string) error
}

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"password"`
	IsEmailVerified bool               `bson:"is_email_verified" json:"is_email_verified"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	FirstName       string             `bson:"first_name" json:"first_name"`
	LastName        string             `bson:"last_name" json:"last_name"`
	Role            string             `bson:"role" json:"role"`
}

type JwtClaims struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Role string `json:"role"`
	jwt.StandardClaims
}
