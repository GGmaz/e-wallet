package model

import "github.com/GGmaz/wallet-arringo/pkg/enums"

type Account struct {
	ID        int64               `json:"id" gorm:"primaryKey;type=varchar(50)"`
	AccNumber string              `json:"acc_number" gorm:"not null"`
	Balance   float64             `json:"balance" gorm:"not null"`
	Status    enums.AccountStatus `json:"status" gorm:"not null"`
	UserId    int64               `json:"user_id" gorm:"not null"`
}

func (tx Account) GetName() string {
	return "Account"
}
