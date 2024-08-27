package controller

import (
	"loan-tracker/domain"
	"loan-tracker/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	// "google.golang.org/api/classroom/v1"
)

type LoanHandler struct {
	loanUsecase *usecase.LoanUsecase
}

func NewLoanHandler(uc *usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{
		loanUsecase: uc,
	}
}

func (h *LoanHandler) ApplyForLoan(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("claim").(domain.JwtClaims)
	userID := claims.UserID


	loan, err := h.loanUsecase.ApplyForLoan(userID.Hex(), req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}

// delivery/loan_handler.go
func (h *LoanHandler) GetLoanStatus(c *gin.Context) {
	loanID := c.Param("id")
	claims := c.MustGet("claim").(domain.JwtClaims)
	userID := claims.UserID
	

	loan, err := h.loanUsecase.GetLoanByID(loanID, userID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	c.JSON(http.StatusOK, loan)
}
func (h *LoanHandler) GetAllLoans(c *gin.Context) {
	status := c.Query("status")
	order := c.Query("order")
	claims := c.MustGet("claim").(domain.JwtClaims)
	
	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	loans, err := h.loanUsecase.GetAllLoans(status, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loans)
}
func (h *LoanHandler) DeleteLoan(c *gin.Context) {
	loanID := c.Param("id")

	claims := c.MustGet("claim").(domain.JwtClaims)
	
	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	err := h.loanUsecase.DeleteLoan(loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}

func (h *LoanHandler) UpdateLoanStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loanID := c.Param("id")
	claims := c.MustGet("claim").(domain.JwtClaims)
	
	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}



	err := h.loanUsecase.UpdateLoanStatus(loanID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan status updated successfully"})
}
