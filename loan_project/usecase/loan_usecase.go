// usecase/loan_usecase.go
package usecase

import (
	"loan-tracker/domain"
	"time"

	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanUsecase struct {
	loanRepo domain.LoanRepository
}

func NewLoanUsecase(repo domain.LoanRepository) *LoanUsecase {
	return &LoanUsecase{
		loanRepo: repo,
	}
}

func (u *LoanUsecase) ApplyForLoan(userID string, amount float64) (*domain.Loan, error) {
	loanID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	loan := &domain.Loan{
		UserID:    loanID,
		Amount:    amount,
		Status:    "pending",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err = u.loanRepo.CreateLoan(loan)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (u *LoanUsecase) GetLoanByID(loanID string, userID string) (*domain.Loan, error) {
	loan, err := u.loanRepo.GetLoanByID(loanID, userID)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

// usecase/loan_usecase.go
func (u *LoanUsecase) GetAllLoans(status, order string) ([]*domain.Loan, error) {
	loans, err := u.loanRepo.GetAllLoans(status, order)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (u *LoanUsecase) UpdateLoanStatus(loanID string, status string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}

	// Update the loan status in the repository
	err := u.loanRepo.UpdateLoanStatus(loanID, status)
	if err != nil {
		return err
	}

	return nil
}
func (u *LoanUsecase) DeleteLoan(loanID string) error {
	return u.loanRepo.DeleteLoan(loanID)
}
