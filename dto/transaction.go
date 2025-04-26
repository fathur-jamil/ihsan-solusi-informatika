package dto

type TransactionRequest struct {
	AccountNumber string  `json:"no_rekening" validate:"required"`
	Amount        float64 `json:"nominal" validate:"required"`
}

type TransactionResponse struct {
	Balance *float64 `json:"saldo,omitempty"`
	Remark  string   `json:"remark,omitempty"`
}
