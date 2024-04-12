package models

type AddMoneyReq struct {
	UserID    int64   `json:"user_id"`
	AccNumber string  `json:"acc_number"`
	Amount    float64 `json:"amount"`
}

type TransferMoneyReq struct {
	From   string  `json:"from_acc_id"`
	To     string  `json:"to_acc_id"`
	Amount float64 `json:"amount"`
}
