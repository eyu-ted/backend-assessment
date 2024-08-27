package routes

import (
	"loan-tracker/delivery/controller"
	"loan-tracker/repository"
	"loan-tracker/usecase" // If you have middleware

	"loan-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)
func SetupUserRouter(db *mongo.Database) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := controller.NewUserHandler(userUsecase)

	
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.GET("/verify-email", userHandler.VerifyEmail)
		userRoutes.POST("/login", userHandler.LoginUser)
		userRoutes.POST("/token/refresh", userHandler.RefreshToken)
		userRoutes.GET("/profile", userHandler.GetUserProfile)
		userRoutes.POST("/request-password-reset", userHandler.PasswordResetRequest)
		userRoutes.POST("/password-update", userHandler.ResetPassword)
	}

	adminRoutes := r.Group("/admin")
	{
		adminRoutes.Use(middleware.AuthMidd)
		adminRoutes.GET("/users", userHandler.ViewAllUsers)
		adminRoutes.DELETE("/users/:id", userHandler.DeleteUser)
	}

	return r
}
