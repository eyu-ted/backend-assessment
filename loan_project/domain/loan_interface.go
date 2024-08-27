package domain

// "go.mongodb.org/mongo-driver/bson/primitive"

type LoanUsecase interface {
	ApplyForLoan(userID string, amount float64) (*Loan, error)
	GetLoanByID(loanID string, userID string) (*Loan, error)
	GetAllLoans(status, order string) ([]*Loan, error)
	UpdateLoanStatus(loanID string, status string) error
	DeleteLoan(loanID string) error
}

type LoanRepository interface {
	CreateLoan(loan *Loan) error
	GetLoanByID(loanID string, userID string) (*Loan, error)
	GetAllLoans(status, order string) ([]*Loan, error)
	UpdateLoanStatus(loanID string, status string) error
	DeleteLoan(loanID string) error
}
