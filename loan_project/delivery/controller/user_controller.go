package controller

import (
	"net/http"
	"time"

	"loan-tracker/domain"
	"loan-tracker/internal/utils"
	"loan-tracker/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(uUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: uUsecase,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUsecase.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.GenerateJWT(&user, "email_verification", 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	emailConfig := utils.EmailConfig{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    587,
		SenderEmail: "eyuted81619@gmail.com",
		SenderName:  "Your Eyoelll",
		SenderPass:  "bdma dfvq tzqw vizv",
	}

	if err := utils.SendVerificationEmail(emailConfig, user.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// delivery/user_handler.go (continued)

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and email are required"})
		return
	}

	if err := h.userUsecase.VerifyEmail(token, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.userUsecase.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.userUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := c.GetString("userID") // Assuming userID is extracted from the token

	user, err := h.userUsecase.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserHandler) PasswordResetRequest(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	token, err := utils.GenerateJWT(&domain.User{Email: req.Email}, "reset", 60*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	emailConfig := utils.EmailConfig{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    587,
		SenderEmail: "eyuted81619@gmail.com",
		SenderName:  "Your Eyoelll",
		SenderPass:  "bdda hfvh tzqq rinv",
	}
	err = utils.SendPasswordResetEmail(emailConfig, req.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully updated"})
}

func (h *UserHandler) ViewAllUsers(c *gin.Context) {
	per := c.MustGet("claim").(domain.JwtClaims)

	if per.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

	per := c.MustGet("claim").(domain.JwtClaims)

	if per.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := c.Param("id")

	err := h.userUsecase.DeleteUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
