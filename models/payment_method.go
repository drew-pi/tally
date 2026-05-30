package models

type PaymentMethodType string

const (
	Credit   PaymentMethodType = "credit"
	Debit    PaymentMethodType = "debit"
	Cheque   PaymentMethodType = "cheque"
	Transfer PaymentMethodType = "transfer"
)

func (p PaymentMethodType) IsValid() bool {
	switch p {
	case Credit, Debit, Cheque, Transfer:
		return true
	}
	return false
}

type PaymentMethod struct {
	ID     int               `json:"id"`
	Type   PaymentMethodType `json:"type"`
	BankID int               `json:"bank_id"`
}
