package models

import "github.com/GGmaz/wallet-arringo/internal/db/model"

type AddMoneyRes struct {
	Balance float64 `json:"updated_balance"`
}

type AccDataRes struct {
	Accounts []model.Account `json:"accounts"`
}
