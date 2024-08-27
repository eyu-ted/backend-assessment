// usecase/user_usecase.go
package usecase

import (
	"errors"
	"time"

	"loan-tracker/domain"
	"loan-tracker/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(uRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: uRepo,
	}
}

func (u *UserUsecase) RegisterUser(user *domain.User) error {
	// Check if user already exists
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set default fields
	// user.ID = generateUUID()
	user.CreatedAt = time.Now()
	user.IsEmailVerified = false
	user.Role = "user"

	// Save user
	return u.userRepo.CreateUser(user)
}

// usecase/user_usecase.go (continued)

func (u *UserUsecase) VerifyEmail(token, email string) error {
	// Get user by email
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// Check if token is valid (This would involve your token verification logic)
	if valid, err := utils.ValidateToken(token, email); err != nil {
		return err
	} else if valid {
		user.IsEmailVerified = true
		return u.userRepo.UpdateUser(user)

	} else {
		return errors.New("invalid token")
		// token is invalid, handle accordingly
	}

}

func (u *UserUsecase) LoginUser(email, password string) (string, string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !user.IsEmailVerified {
		return "", "", errors.New("email not verified")
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateJWT(user, "access", 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateJWT(user, "refresh", 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *UserUsecase) RefreshToken(refreshToken string) (string, error) {
	// Validate and parse the refresh token
	claims, err := utils.ParseJWT(refreshToken)
	if err != nil {
		return "", errors.New("invalid token")
	}

	if claims.Type != "refresh" {
		return "", errors.New("invalid token type")
	}

	// Generate a new access token
	user, err := u.userRepo.GetUserByEmail(claims.Email)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateJWT(user, "access", 15*time.Minute)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u *UserUsecase) GetUserProfile(userID string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(userID) // Implement GetUserByID in the repository
	if err != nil {
		return nil, err
	}
	return user, nil
}

// usecase/user_usecase.go
func (u *UserUsecase) ResetPassword(token, newPassword string) error {
	claims, err := utils.ParseJWT(token) // Implement token verification
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetUserByEmail(claims.Email)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) DeleteUserByID(userID string) error {
	err := u.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserUsecase) GetAllUsers() ([]*domain.User, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
