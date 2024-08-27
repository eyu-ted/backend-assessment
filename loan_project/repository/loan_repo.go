package repository

import (
	"context"
	"loan-tracker/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoLoanRepository struct {
	db *mongo.Collection
}

func NewMongoLoanRepository(db *mongo.Database) domain.LoanRepository {
	return &mongoLoanRepository{
		db: db.Collection("loans"),
	}
}

func (r *mongoLoanRepository) CreateLoan(loan *domain.Loan) error {
	loan.ID = primitive.NewObjectID()
	_, err := r.db.InsertOne(context.TODO(), loan)
	return err
}

func (r *mongoLoanRepository) GetLoanByID(loanID string, userID string) (*domain.Loan, error) {
	var loan domain.Loan
	objID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return nil, err
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "user_id": userObjID}
	err = r.db.FindOne(context.TODO(), filter).Decode(&loan)
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

// repository/loan_repository_mongo.go
func (r *mongoLoanRepository) GetAllLoans(status, order string) ([]*domain.Loan, error) {
	var loans []*domain.Loan
	filter := bson.M{}
	if status != "" && status != "all" {
		filter["status"] = status
	}

	sort := bson.D{{"created_at", 1}} // Default to ascending order
	if order == "desc" {
		sort = bson.D{{"created_at", -1}}
	}

	cursor, err := r.db.Find(context.TODO(), filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var loan domain.Loan
		if err = cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

// repository/loan_repository_mongo.go
func (r *mongoLoanRepository) DeleteLoan(loanID string) error {
	objID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = r.db.DeleteOne(context.TODO(), filter)
	return err
}

func (r *mongoLoanRepository) UpdateLoanStatus(loanID string, status string) error {
	objID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now().Unix()}}

	_, err = r.db.UpdateOne(context.TODO(), filter, update)
	return err
}
