package utils

import (
	"errors"
	"fmt"
	"loan-tracker/domain"
	"time"

	// jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = "your-secret-key"

func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateJWT(user *domain.User, tokenType string, duration time.Duration) (string, error) {
	exp := time.Now().Add(time.Hour * duration).Unix()
	claims := domain.JwtClaims{
		Email:  user.Email,
		Type:   tokenType,
		Role:   user.Role,
		UserID: user.ID, // Convert ObjectID to string
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenStr, email string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return false, err
	}
	fmt.Println(token)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["Email"] == email {
			return true, nil
		}
	}

	return false, errors.New("invalid token")
}

func ParseJWT(token string) (*domain.JwtClaims, error) {
	claims := &domain.JwtClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil

}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
