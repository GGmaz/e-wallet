package model

import (
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"time"
)

type Transaction struct {
	ID              int64          `json:"id" gorm:"primaryKey, type=varchar(50)"`
	UserID          int64          `gorm:"not null" json:"user_id"`
	AccountNumber   string         `gorm:"not null" json:"account_number"`
	Amount          float64        `gorm:"not null" json:"amount"`
	Balance         float64        `gorm:"not null" json:"balance"`
	TransactionType enums.TxType   `gorm:"not null" json:"transaction_type"`
	Status          enums.TxStatus `gorm:"not null" json:"status"`
	CreatedAt       time.Time      `json:"created_at" gorm:"default:now()"`
	PairedAccNum    string         `json:"paired_tx_num"`
}

func (tx Transaction) GetName() string {
	return "Transaction"
}
