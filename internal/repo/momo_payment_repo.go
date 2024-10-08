package repo

type MoMoRepo interface {
	CreatePayment(orderID, amount string) (string, error)
	CheckPaymentStatus(orderID string) (bool, error)
	RefundPayment(orderID, amount string) (string, error)
}

type momoRepo struct {
	endpoint    string
	partnerCode string
	accessKey   string
	secrectKey  string
}
