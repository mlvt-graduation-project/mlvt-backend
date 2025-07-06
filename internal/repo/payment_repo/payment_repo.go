package payment_repo

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.PaymentTransaction) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.PaymentTransaction, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*entity.PaymentTransaction, error)
	GetByUserID(ctx context.Context, userID uint64) ([]entity.PaymentTransaction, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status entity.PaymentStatus) error
	MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error
	GetPendingPayments(ctx context.Context) ([]entity.PaymentTransaction, error)
}

type paymentRepo struct {
	adapter *mongodb.MongoDBAdapter[entity.PaymentTransaction]
}

func NewPaymentRepo(mongoClient *mongodb.MongoDBClient) PaymentRepository {
	return &paymentRepo{
		adapter: mongodb.NewMongoDBAdapter[entity.PaymentTransaction](
			mongoClient.GetClient(),
			"mlvt",
			"payment_transactions",
		),
	}
}

func (r *paymentRepo) Create(ctx context.Context, payment *entity.PaymentTransaction) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	insertedID, err := r.adapter.InsertOne(*payment)
	if err != nil {
		return err
	}

	payment.ID = insertedID
	return nil
}

func (r *paymentRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.PaymentTransaction, error) {
	filter := bson.M{"_id": id}
	return r.adapter.FindOne(filter)
}

func (r *paymentRepo) GetByTransactionID(ctx context.Context, transactionID string) (*entity.PaymentTransaction, error) {
	filter := bson.M{"transaction_id": transactionID}
	return r.adapter.FindOne(filter)
}

func (r *paymentRepo) GetByUserID(ctx context.Context, userID uint64) ([]entity.PaymentTransaction, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	return r.adapter.Find(filter, opts)
}

func (r *paymentRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status entity.PaymentStatus) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}

	return r.adapter.UpdateOne(filter, update)
}

func (r *paymentRepo) MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"status":       entity.PaymentStatusCompleted,
		"completed_at": time.Now(),
		"updated_at":   time.Now(),
	}

	return r.adapter.UpdateOne(filter, update)
}

func (r *paymentRepo) GetPendingPayments(ctx context.Context) ([]entity.PaymentTransaction, error) {
	filter := bson.M{"status": entity.PaymentStatusPending}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	return r.adapter.Find(filter, opts)
}
