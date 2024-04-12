package repo

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"gorm.io/gorm"
)

func GetUserByMail(db *gorm.DB, mail string) (*model.User, error) {
	user := &model.User{Email: mail}
	res := db.Preload("Accounts").First(&user, "email = ?", mail)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func HasUserAccount(db *gorm.DB, userId int64, accNum string) bool {
	user := &model.User{}
	acc := &model.Account{}
	db.First(&user, userId)
	db.First(&acc, "acc_number = ?", accNum)
	return acc.UserId == user.ID
}
