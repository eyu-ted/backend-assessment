package routes

import (
	"loan-tracker/delivery/controller"
	"loan-tracker/internal/middleware"
	"loan-tracker/repository"
	"loan-tracker/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupLoanRouter(db *mongo.Database, r *gin.Engine) *gin.Engine {

	loanRepo := repository.NewMongoLoanRepository(db)

	loanUsecase := usecase.NewLoanUsecase(loanRepo)
	loanHandler := controller.NewLoanHandler(loanUsecase)

	// Loan routes for regular users
	loanRoutes := r.Group("/loans")
	{
		loanRoutes.Use(middleware.AuthMidd)
		loanRoutes.POST("/", loanHandler.ApplyForLoan)
		loanRoutes.GET("/:id", loanHandler.GetLoanStatus)
	}

	// Loan routes for admins
	adminRoutes := r.Group("/admin")
	{
		adminRoutes.Use(middleware.AuthMidd)
		adminRoutes.GET("/loans", loanHandler.GetAllLoans)
		adminRoutes.PATCH("/loans/:id/status", loanHandler.UpdateLoanStatus)
		adminRoutes.DELETE("/loans/:id", loanHandler.DeleteLoan)
	}

	return r
}
