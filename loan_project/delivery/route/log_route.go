package routes

import (
	"loan-tracker/delivery/controller"
	"loan-tracker/internal/middleware"
	"loan-tracker/repository"
	"loan-tracker/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupLogRouter(db *mongo.Database, r *gin.Engine) *gin.Engine {

	logRepo := repository.NewLogRepository(db)
	logUsecase := usecase.NewLogUsecase(logRepo)
	logHandler := controller.NewLogHandler(logUsecase)

	logRoutes := r.Group("/logs")
	{
		logRoutes.Use(middleware.AuthMidd)
		logRoutes.GET("/", logHandler.GetLogs)
		logRoutes.POST("/", logHandler.CreateLog)
	}

	return r
}
