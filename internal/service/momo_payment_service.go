package service

type PaymentService interface {
	GeneratePaymentQRCode(orderID, amount string) ([]byte, error)
	CheckPaymentStatus(orderID string) (bool, error)
	RefundPayment(orderID, amount string) (string, error)
}
