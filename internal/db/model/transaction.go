package model

import (
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"time"
)

type Transaction struct {
	ID              int64        `json:"id" gorm:"primaryKey, type=varchar(50)"`
	UserID          int64        `gorm:"not null" json:"user_id"`
	Amount          float64      `gorm:"not null" json:"amount"`
	Balance         float64      `gorm:"not null" json:"balance"`
	TransactionType enums.TxType `gorm:"not null" json:"transaction_type"`
	CreatedAt       time.Time    `json:"created_at" gorm:"default:now()"`
}

func (tx Transaction) GetName() string {
	return "Transaction"
}
