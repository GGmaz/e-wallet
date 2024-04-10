package repo

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"gorm.io/gorm"
)

func GetBalanceForUserMail(db *gorm.DB, mail string) (float64, error) {
	user := &model.User{Email: mail}
	res := db.First(&user, "email = ?", mail)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.Balance, nil
}

func Create(gormDB *gorm.DB, t *model.User) *gorm.DB {
	return gormDB.Create(&t)
}
