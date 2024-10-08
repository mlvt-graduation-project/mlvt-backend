package entity

// TransactionLog represents a log entry for a transaction
type TransactionLog struct {
	OrderID       string
	PaymentMethod string
	Action        string
	Status        string
	Details       string
}
