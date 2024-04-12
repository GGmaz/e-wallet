package repo

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAccountByNum(db *gorm.DB, acc *model.Account, accNum string, preload ...string) *gorm.DB {
	for _, m := range preload {
		db = db.Preload(m)
	}
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc, "acc_number = ?", accNum)
}
