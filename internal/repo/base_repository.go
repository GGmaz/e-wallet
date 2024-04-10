package repo

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type Repo[T model.BaseInterface] struct {
}

func (r Repo[T]) Create(db *gorm.DB, t *T) *gorm.DB {
	log.Printf("%v", t)
	return db.Create(&t)
}

func (r Repo[T]) GetById(db *gorm.DB, t *T, id string, preload ...string) *gorm.DB {
	for _, m := range preload {
		db = db.Preload(m)
	}
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t, "id = ?", id)
}

func (r Repo[T]) Save(db *gorm.DB, t *T) *gorm.DB {
	return db.Clauses(clause.Returning{}).Save(&t)
}

func (r Repo[T]) Update(db *gorm.DB, t *T, id int64) *gorm.DB {
	return db.Model(t).Clauses(clause.Returning{}).Where("id = ?", id).Updates(t)
}
