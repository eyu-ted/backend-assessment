package controller

import (
	"loan-tracker/domain"
	"loan-tracker/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logUsecase *usecase.LogUsecase
}

func NewLogHandler(logUsecase *usecase.LogUsecase) *LogHandler {
	return &LogHandler{
		logUsecase: logUsecase,
	}
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	logs, err := h.logUsecase.GetLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
func (h *LogHandler) CreateLog(c *gin.Context) {
	var log domain.Log
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.logUsecase.CreateLog(log.EventType, log.UserID.Hex(), log.Details, log.Success); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "log created successfully"})
}
